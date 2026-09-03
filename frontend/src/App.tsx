import { BrowserRouter, Routes, Route } from 'react-router-dom'
import Layout from './components/Layout'
import Audits from './pages/Audits'
import AuditDetail from './pages/AuditDetail'
import NewAudit from './pages/NewAudit'

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Layout><Audits /></Layout>} />
        <Route path="/audits/:id" element={<Layout><AuditDetail /></Layout>} />
        <Route path="/new" element={<Layout><NewAudit /></Layout>} />
      </Routes>
    </BrowserRouter>
  )
}
