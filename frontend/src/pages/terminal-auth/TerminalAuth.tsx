import { useState } from 'react'
import { api } from '../../services/api'

const TerminalAuth = () => {
  const [cardNumber, setCardNumber] = useState('')
  const [price, setPrice] = useState('')
  const [terminalSerial, setTerminalSerial] = useState('')
  const [result, setResult] = useState<{ authorized: boolean; message: string } | null>(null)
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setResult(null)

    try {
      const response = await api('/api/v1/terminal/auth', {
        method: 'POST',
        body: JSON.stringify({
          card_number: cardNumber,
          price: parseFloat(price),
          terminal_serial_number: terminalSerial
        })
      })
      setResult(response)
    } catch (err: any) {
      setResult({
        authorized: false,
        message: err.message || 'Authorization failed'
      })
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="max-w-lg mx-auto">
      <div className="bg-white rounded-2xl shadow-xl p-8">
        <h2 className="text-3xl font-bold text-center mb-8">💳 Terminal Authorization</h2>

        <form onSubmit={handleSubmit} className="space-y-6">
          <div>
            <label className="block text-sm font-medium mb-2">Card Number</label>
            <input
              type="text"
              value={cardNumber}
              onChange={(e) => setCardNumber(e.target.value)}
              className="w-full px-4 py-3 border rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="12345678901234"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium mb-2">Amount</label>
            <input
              type="number"
              value={price}
              onChange={(e) => setPrice(e.target.value)}
              className="w-full px-4 py-3 border rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="150.00"
              step="0.01"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium mb-2">Terminal Serial Number (optional)</label>
            <input
              type="text"
              value={terminalSerial}
              onChange={(e) => setTerminalSerial(e.target.value)}
              className="w-full px-4 py-3 border rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="TERM001"
            />
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full bg-gradient-to-r from-blue-600 to-indigo-600 text-white py-4 rounded-xl text-lg font-semibold hover:from-blue-700 hover:to-indigo-700 disabled:opacity-50"
          >
            {loading ? 'Processing...' : 'Authorize Payment'}
          </button>
        </form>

        {result && (
          <div className={`mt-8 p-6 rounded-xl text-center text-lg font-medium ${
            result.authorized ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
          }`}>
            {result.authorized ? '✅' : '❌'} {result.message}
          </div>
        )}
      </div>
    </div>
  )
}

export default TerminalAuth