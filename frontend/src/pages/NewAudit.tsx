export default function NewAudit() {
  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-gray-900">New Audit</h1>
      <div className="bg-white rounded-lg shadow p-6">
        <form className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700">Website URL</label>
            <input type="url" className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500" placeholder="https://example.com" />
          </div>
          <button type="submit" className="inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700">
            Run Audit
          </button>
        </form>
      </div>
    </div>
  )
}
