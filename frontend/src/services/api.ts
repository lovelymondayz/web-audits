import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  headers: { 'Content-Type': 'application/json' },
})

export async function getAudits() {
  const res = await api.get('/audits')
  return res.data
}

export async function getAudit(id: string) {
  const res = await api.get(`/audits/${id}`)
  return res.data
}

export async function createAudit(url: string) {
  const res = await api.post('/audits', { url })
  return res.data
}
