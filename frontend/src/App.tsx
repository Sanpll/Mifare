import { BrowserRouter as Router } from 'react-router-dom'
import Navbar from './components/Layout/Navbar'
import AppRouter from './router'

function App() {
  return (
    <Router>
      <div className="min-h-screen bg-gray-50">
        <Navbar />
        <main className="max-w-7xl mx-auto px-4 py-6">
          <AppRouter />
        </main>
      </div>
    </Router>
  )
}

export default App