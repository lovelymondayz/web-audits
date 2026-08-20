
from fastapi import FastAPI, HTTPException, BackgroundTasks
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
import sqlite3, os, uuid, json, requests
from datetime import datetime
from pathlib import Path

app = FastAPI(title="Web Audits API", version="1.0.0")
app.add_middleware(CORSMiddleware, allow_origins=["*"], allow_methods=["*"], allow_headers=["*"])

DATABASE_PATH = os.getenv("DATABASE_PATH", "/app/data/audits.db")
Path(DATABASE_PATH).parent.mkdir(parents=True, exist_ok=True)

def init_db():
    conn = sqlite3.connect(DATABASE_PATH)
    c = conn.cursor()
    c.execute("CREATE TABLE IF NOT EXISTS audits (id TEXT PRIMARY KEY, url TEXT, status TEXT DEFAULT 'pending', score INTEGER, created_at TEXT, completed_at TEXT)")
    c.execute("CREATE TABLE IF NOT EXISTS audit_categories (id TEXT PRIMARY KEY, audit_id TEXT, category TEXT, score INTEGER, findings TEXT)")
    c.execute("CREATE TABLE IF NOT EXISTS audit_recommendations (id TEXT PRIMARY KEY, audit_id TEXT, priority TEXT, title TEXT, description TEXT)")
    conn.commit(); conn.close()

def get_db():
    conn = sqlite3.connect(DATABASE_PATH); conn.row_factory = sqlite3.Row; return conn

class AuditReq(BaseModel):
    url: str

@app.get("/health")
async def health(): return {"status": "healthy"}

@app.get("/")
async def root(): return {"service": "Web Audits API", "version": "1.0.0"}

def run_audit(aid, url):
    conn = get_db(); c = conn.cursor()
    try:
        resp = requests.get(url, timeout=15, headers={"User-Agent":"Mozilla/5.0"})
        html = resp.text
        findings = {}
        # Performance: count resources
        scripts = html.count("<script"); styles = html.count("<link"); imgs = html.count("<img")
        perf_score = max(0, 100 - scripts*2 - styles*3 - imgs)
        findings["performance"] = {"score": min(100, perf_score), "findings": f"Found {scripts} scripts, {styles} stylesheets, {imgs} images"}

        # Mobile: viewport meta
        mobile_score = 100 if 'viewport' in html.lower() else 40
        findings["mobile"] = {"score": mobile_score, "findings": "Viewport meta tag " + ("found" if mobile_score == 100 else "MISSING")}

        # SEO: title, meta desc, h1
        has_title = "<title>" in html.lower()
        has_meta_desc = 'name="description"' in html.lower()
        has_h1 = "<h1" in html.lower()
        seo_score = (has_title*30 + has_meta_desc*35 + has_h1*35)
        findings["seo"] = {"score": seo_score, "findings": f"Title: {'yes' if has_title else 'no'}, Meta desc: {'yes' if has_meta_desc else 'no'}, H1: {'yes' if has_h1 else 'no'}"}

        # CTA: buttons, forms
        buttons = html.count("<button") + html.count('class="btn"')
        forms = html.count("<form")
        cta_score = min(100, buttons*10 + forms*20 + 50)
        findings["cta"] = {"score": cta_score, "findings": f"Found {buttons} buttons, {forms} forms"}

        # Accessibility: alt tags
        alts = html.count('alt="')
        no_alts = html.count('<img') - alts
        acc_score = max(0, 100 - no_alts*5)
        findings["accessibility"] = {"score": acc_score, "findings": f"Images with alt: {alts}, without: {no_alts}"}

        total = int(perf_score*0.25 + mobile_score*0.2 + seo_score*0.2 + cta_score*0.2 + acc_score*0.15)
        c.execute("UPDATE audits SET status=?,score=?,completed_at=? WHERE id=?", ("completed", total, datetime.utcnow().isoformat(), aid))

        for cat, data in findings.items():
            cid = str(uuid.uuid4())
            c.execute("INSERT INTO audit_categories (id,audit_id,category,score,findings) VALUES (?,?,?,?,?)",
                (cid, aid, cat, data["score"], data["findings"]))

        recs = []
        if mobile_score < 80: recs.append(("high","Add viewport meta tag","Mobile usability is critical for SEO"))
        if seo_score < 70: recs.append(("high","Fix SEO basics","Add title, meta description, and H1 tags"))
        if perf_score < 60: recs.append(("medium","Optimize performance","Reduce scripts and images"))
        if cta_score < 60: recs.append(("medium","Add clear CTAs","Include more buttons and forms"))
        if acc_score < 80: recs.append(("low","Add alt text to images","Improve accessibility"))
        for pri, title, desc in recs:
            rid = str(uuid.uuid4())
            c.execute("INSERT INTO audit_recommendations (id,audit_id,priority,title,description) VALUES (?,?,?,?,?)",
                (rid, aid, pri, title, desc))

        conn.commit()
    except Exception as e:
        c.execute("UPDATE audits SET status=?,completed_at=? WHERE id=?", ("failed", datetime.utcnow().isoformat(), aid))
        conn.commit()
    finally:
        conn.close()

@app.post("/audits")
async def create_audit(req: AuditReq, bg: BackgroundTasks):
    aid = str(uuid.uuid4())
    conn = get_db()
    c = conn.cursor()
    c.execute("INSERT INTO audits (id,url,status,created_at) VALUES (?,?,?,?)", (aid, req.url, "running", datetime.utcnow().isoformat()))
    conn.commit(); conn.close()
    bg.add_task(run_audit, aid, req.url)
    return {"audit_id": aid, "status": "running"}

@app.get("/audits")
async def list_audits():
    conn = get_db()
    audits = [dict(r) for r in conn.execute("SELECT * FROM audits ORDER BY created_at DESC").fetchall()]
    conn.close(); return {"audits": audits}

@app.get("/audits/{aid}")
async def get_audit(aid: str):
    conn = get_db()
    c = conn.cursor()
    c.execute("SELECT * FROM audits WHERE id=?", (aid,))
    audit = c.fetchone()
    if not audit: conn.close(); raise HTTPException(404, "Audit not found")
    c.execute("SELECT * FROM audit_categories WHERE audit_id=?", (aid,))
    categories = [dict(r) for r in c.fetchall()]
    c.execute("SELECT * FROM audit_recommendations WHERE audit_id=? ORDER BY priority", (aid,))
    recs = [dict(r) for r in c.fetchall()]
    conn.close()
    result = dict(audit); result["categories"] = categories; result["recommendations"] = recs
    return result

@app.on_event("startup")
async def startup(): init_db()
