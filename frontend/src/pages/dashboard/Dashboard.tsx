import { useAuth } from '../../context/AuthContext'

const Dashboard = () => {
  const { user } = useAuth()

  return (
    <div className="text-center">
      <h1 className="text-4xl font-bold mb-4">Welcome to Mifare System</h1>
      <p className="text-xl text-gray-600 mb-8">
        Hello, <span className="font-semibold">{user?.username}</span>!
      </p>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 max-w-4xl mx-auto">
        <div className="bg-white p-6 rounded-xl shadow">
          <h3 className="text-lg font-semibold mb-2">Cards</h3>
          <p className="text-3xl font-bold text-blue-600">12</p>
          <p className="text-sm text-gray-500">Active cards</p>
        </div>

        <div className="bg-white p-6 rounded-xl shadow">
          <h3 className="text-lg font-semibold mb-2">Transactions</h3>
          <p className="text-3xl font-bold text-green-600">87</p>
          <p className="text-sm text-gray-500">This month</p>
        </div>

        <div className="bg-white p-6 rounded-xl shadow">
          <h3 className="text-lg font-semibold mb-2">Terminals</h3>
          <p className="text-3xl font-bold text-purple-600">5</p>
          <p className="text-sm text-gray-500">Connected</p>
        </div>
      </div>

      <div className="mt-12">
        <a href="/terminal/auth" 
           className="inline-block bg-blue-600 text-white px-8 py-4 rounded-xl text-lg font-medium hover:bg-blue-700">
          Go to Terminal Authorization →
        </a>
      </div>
    </div>
  )
}

export default Dashboard