import React, { useState, useEffect } from 'react'

// ---- helper: декодировать JWT (без валидации) ----
const decodeJWT = (token) => {
    try {
        const base64Url = token.split('.')[1]
        const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/')
        const jsonPayload = decodeURIComponent(atob(base64).split('').map(c => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2)).join(''))
        return JSON.parse(jsonPayload)
    } catch (e) {
        return null
    }
}

let authToken = localStorage.getItem('token') || null
let isAdmin = false
if (authToken) {
    const payload = decodeJWT(authToken)
    isAdmin = payload?.is_admin === true
}

const request = async (endpoint, options = {}) => {
    const headers = {
        'Content-Type': 'application/json',
        ...options.headers
    }
    if (authToken) {
        headers['Authorization'] = `Bearer ${authToken}`
    }
    const res = await fetch(endpoint, { ...options, headers })
    if (res.status === 401 && authToken) {
        localStorage.removeItem('token')
        authToken = null
        isAdmin = false
        window.location.reload()
    }
    return res
}

// ---- Компонент логина с проверкой админа ----
const Login = ({ onLogin }) => {
    const [username, setUsername] = useState('')
    const [password, setPassword] = useState('')
    const [error, setError] = useState('')

    const handleSubmit = async (e) => {
        e.preventDefault()
        const res = await fetch('/auth/sign-in', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password })
        })
        if (res.ok) {
            const data = await res.json()
            const token = data.token
            const payload = decodeJWT(token)
            if (payload && payload.is_admin === true) {
                authToken = token
                isAdmin = true
                localStorage.setItem('token', authToken)
                onLogin()
            } else {
                setError('Доступ разрешён только администраторам')
            }
        } else {
            setError('Неверное имя пользователя или пароль')
        }
    }

    return (
        <div className="min-h-screen flex items-center justify-center bg-gray-50">
            <div className="bg-white p-8 rounded shadow-md w-96">
                <h1 className="text-2xl font-bold mb-6">Вход в админ-панель</h1>
                <form onSubmit={handleSubmit}>
                    <input type="text" placeholder="Имя пользователя" value={username} onChange={e => setUsername(e.target.value)} className="border p-2 w-full mb-4 rounded" required />
                    <input type="password" placeholder="Пароль" value={password} onChange={e => setPassword(e.target.value)} className="border p-2 w-full mb-4 rounded" required />
                    {error && <p className="text-red-500 mb-4">{error}</p>}
                    <button type="submit" className="bg-blue-600 text-white py-2 px-4 rounded w-full">Войти</button>
                </form>
            </div>
        </div>
    )
}

// ---- Универсальная CRUD таблица (для простых сущностей) ----
const CrudTable = ({ title, apiPath, fields, getId = (item) => item.id }) => {
    const [items, setItems] = useState([])
    const [form, setForm] = useState({})
    const [editingId, setEditingId] = useState(null)
    const [loading, setLoading] = useState(true)

    const fetchItems = async () => {
        const res = await request(`${apiPath}`)
        if (res.ok) {
            const data = await res.json()
            // ключи ожидаем в data.keys, пользователи в data.users, терминалы в data.terminals, транзакции в data.transactions
            const list = data[Object.keys(data)[0]] || data
            setItems(list)
        }
        setLoading(false)
    }

    useEffect(() => { fetchItems() }, [apiPath])

    const handleSubmit = async (e) => {
        e.preventDefault()
        const method = editingId ? 'PUT' : 'POST'
        const url = editingId ? `${apiPath}/${editingId}` : apiPath
        const res = await request(url, { method, body: JSON.stringify(form) })
        if (res.ok) {
            fetchItems()
            setForm({})
            setEditingId(null)
        } else {
            const err = await res.text()
            alert('Ошибка: ' + err)
        }
    }

    const handleDelete = async (id) => {
        if (confirm('Удалить?')) {
            const res = await request(`${apiPath}/${id}`, { method: 'DELETE' })
            if (res.ok) fetchItems()
            else alert('Ошибка удаления')
        }
    }

    const editItem = (item) => {
        setForm(item)
        setEditingId(getId(item))
    }

    if (loading) return <div>Загрузка...</div>

    return (
        <div className="p-4">
            <h2 className="text-xl font-bold mb-4">{title}</h2>
            <form onSubmit={handleSubmit} className="mb-6 flex flex-wrap gap-2">
                {fields.map(f => (
                    <input key={f.name} type={f.type || 'text'} placeholder={f.label} value={form[f.name] || ''} onChange={e => setForm({ ...form, [f.name]: e.target.value })} className="border p-2 rounded" />
                ))}
                <button type="submit" className="bg-green-600 text-white px-4 py-2 rounded">{editingId ? 'Обновить' : 'Создать'}</button>
                {editingId && <button type="button" onClick={() => { setEditingId(null); setForm({}) }} className="bg-gray-400 text-white px-4 py-2 rounded">Отмена</button>}
            </form>
            <table className="w-full border-collapse border">
                <thead><tr className="bg-gray-100">{fields.map(f => <th key={f.name} className="border p-2">{f.label}</th>)}<th className="border p-2">Действия</th></tr></thead>
                <tbody>
                    {items.map(item => (
                        <tr key={getId(item)}>
                            {fields.map(f => <td key={f.name} className="border p-2">{item[f.name]?.toString()}</td>)}
                            <td className="border p-2">
                                <button onClick={() => editItem(item)} className="bg-yellow-500 text-white px-2 py-1 rounded mr-2">Ред.</button>
                                <button onClick={() => handleDelete(getId(item))} className="bg-red-500 text-white px-2 py-1 rounded">Уд.</button>
                            </td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    )
}

// ---- Таблица для карт (особая логика: owner_name не в создании, is_blocked чекбокс) ----
const CardsTable = () => {
    const [items, setItems] = useState([])
    const [form, setForm] = useState({})
    const [editingId, setEditingId] = useState(null)
    const [loading, setLoading] = useState(true)

    const fetchItems = async () => {
        const res = await request('/api/v1/cards')
        if (res.ok) {
            const data = await res.json()
            setItems(data.cards || [])
        }
        setLoading(false)
    }

    useEffect(() => { fetchItems() }, [])

    const handleSubmit = async (e) => {
        e.preventDefault()
        if (!editingId) {
            const { owner_name, ...createData } = form
            const res = await request('/api/v1/cards', { method: 'POST', body: JSON.stringify(createData) })
            if (res.ok) {
                fetchItems()
                setForm({})
            } else {
                alert('Ошибка создания')
            }
        } else {
            const res = await request(`/api/v1/cards/${editingId}`, { method: 'PUT', body: JSON.stringify(form) })
            if (res.ok) {
                fetchItems()
                setForm({})
                setEditingId(null)
            } else {
                alert('Ошибка обновления')
            }
        }
    }

    const handleDelete = async (id) => {
        if (confirm('Удалить карту?')) {
            const res = await request(`/api/v1/cards/${id}`, { method: 'DELETE' })
            if (res.ok) fetchItems()
            else alert('Ошибка удаления')
        }
    }

    const editItem = (item) => {
        setForm(item)
        setEditingId(item.id)
    }

    if (loading) return <div>Загрузка...</div>

    return (
        <div className="p-4">
            <h2 className="text-xl font-bold mb-4">Карты</h2>
            <form onSubmit={handleSubmit} className="mb-6 flex flex-wrap gap-2">
                <input type="text" placeholder="Номер карты (hex)" value={form.number || ''} onChange={e => setForm({ ...form, number: e.target.value })} className="border p-2 rounded" required={!editingId} />
                <input type="number" step="0.01" placeholder="Баланс" value={form.balance || ''} onChange={e => setForm({ ...form, balance: parseFloat(e.target.value) })} className="border p-2 rounded" required={!editingId} />
                <label className="flex items-center gap-1"><input type="checkbox" checked={form.is_blocked || false} onChange={e => setForm({ ...form, is_blocked: e.target.checked })} /> Заблокирована</label>
                <input type="text" placeholder="Ключ (12 hex)" value={form.key_value || ''} onChange={e => setForm({ ...form, key_value: e.target.value })} className="border p-2 rounded" required={!editingId} />
                {editingId && <input type="text" placeholder="Владелец" value={form.owner_name || ''} onChange={e => setForm({ ...form, owner_name: e.target.value })} className="border p-2 rounded" />}
                <button type="submit" className="bg-green-600 text-white px-4 py-2 rounded">{editingId ? 'Обновить' : 'Создать'}</button>
                {editingId && <button type="button" onClick={() => { setEditingId(null); setForm({}) }} className="bg-gray-400 text-white px-4 py-2 rounded">Отмена</button>}
            </form>
            <table className="w-full border-collapse border">
                <thead><tr className="bg-gray-100"><th className="border p-2">ID</th><th className="border p-2">Номер</th><th className="border p-2">Баланс</th><th className="border p-2">Заблокирована</th><th className="border p-2">Владелец</th><th className="border p-2">Ключ</th><th className="border p-2">Действия</th></tr></thead>
                <tbody>
                    {items.map(card => (
                        <tr key={card.id}>
                            <td className="border p-2">{card.id}</td>
                            <td className="border p-2">{card.number}</td>
                            <td className="border p-2">{card.balance}</td>
                            <td className="border p-2">{card.is_blocked ? 'Да' : 'Нет'}</td>
                            <td className="border p-2">{card.owner_name}</td>
                            <td className="border p-2">{card.key_value}</td>
                            <td className="border p-2">
                                <button onClick={() => editItem(card)} className="bg-yellow-500 text-white px-2 py-1 rounded mr-2">Ред.</button>
                                <button onClick={() => handleDelete(card.id)} className="bg-red-500 text-white px-2 py-1 rounded">Уд.</button>
                            </td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    )
}

// ---- Проверка авторизации транзакции ----
const TransactionAuth = () => {
    const [cardNumber, setCardNumber] = useState('')
    const [price, setPrice] = useState('')
    const [terminalSN, setTerminalSN] = useState('')
    const [result, setResult] = useState(null)

    const checkAuth = async () => {
        const res = await request('/api/v1/terminal/auth', {
            method: 'POST',
            body: JSON.stringify({ card_number: cardNumber, price: parseFloat(price), terminal_serial_number: terminalSN })
        })
        const data = await res.json()
        setResult(data)
    }

    return (
        <div className="p-4">
            <h2 className="text-xl font-bold mb-4">Авторизация транзакции</h2>
            <div className="mb-4 flex gap-2 flex-wrap">
                <input placeholder="Номер карты" value={cardNumber} onChange={e => setCardNumber(e.target.value)} className="border p-2 rounded" />
                <input placeholder="Сумма" value={price} onChange={e => setPrice(e.target.value)} className="border p-2 rounded" />
                <input placeholder="Серийный номер терминала" value={terminalSN} onChange={e => setTerminalSN(e.target.value)} className="border p-2 rounded" />
                <button onClick={checkAuth} className="bg-blue-600 text-white px-4 py-2 rounded">Проверить</button>
            </div>
            {result && (
                <div className={`p-3 rounded ${result.authorized ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}`}>
                    {result.message}
                </div>
            )}
        </div>
    )
}

// ---- Админ-панель ----
const AdminPanel = () => {
    const [activeTab, setActiveTab] = useState('users')

    if (!isAdmin) {
        return <div className="p-8 text-red-600 text-center">Доступ запрещён. Только для администраторов.</div>
    }

    const tabs = [
        { id: 'users', label: 'Пользователи', component: () => <CrudTable title="Пользователи" apiPath="/api/v1/users" fields={[{ name: 'username', label: 'Имя' }]} getId={item => item.id} /> },
        { id: 'keys', label: 'Ключи', component: () => <CrudTable title="Ключи" apiPath="/api/v1/keys" fields={[{ name: 'value', label: 'Значение' }, { name: 'type', label: 'Тип' }, { name: 'description', label: 'Описание' }]} getId={item => item.id} /> },
        { id: 'cards', label: 'Карты', component: CardsTable },
        { id: 'terminals', label: 'Терминалы', component: () => <CrudTable title="Терминалы" apiPath="/api/v1/terminals" fields={[{ name: 'serial_number', label: 'Серийный номер' }, { name: 'address', label: 'Адрес' }, { name: 'name', label: 'Название' }]} getId={item => item.id} /> },
        { id: 'transactions', label: 'Транзакции', component: () => <CrudTable title="Транзакции" apiPath="/api/v1/transactions" fields={[{ name: 'card_number', label: 'Номер карты' }, { name: 'price', label: 'Сумма' }, { name: 'terminal_serial_number', label: 'Терминал' }]} getId={item => item.id} /> },
        { id: 'auth', label: 'Проверка транзакции', component: TransactionAuth }
    ]

    const logout = () => {
        localStorage.removeItem('token')
        authToken = null
        isAdmin = false
        window.location.reload()
    }

    const ActiveComponent = tabs.find(t => t.id === activeTab)?.component || tabs[0].component

    return (
        <div>
            <div className="bg-gray-800 text-white p-4 flex justify-between">
                <div className="flex gap-4 flex-wrap">
                    {tabs.map(tab => (
                        <button key={tab.id} onClick={() => setActiveTab(tab.id)} className={`px-3 py-1 rounded ${activeTab === tab.id ? 'bg-blue-600' : 'hover:bg-gray-700'}`}>{tab.label}</button>
                    ))}
                </div>
                <button onClick={logout} className="bg-red-600 px-3 py-1 rounded">Выйти</button>
            </div>
            <div className="p-4"><ActiveComponent /></div>
        </div>
    )
}

// ---- Точка входа ----
const App = () => {
    const [isLoggedIn, setIsLoggedIn] = useState(!!authToken && isAdmin)
    if (!isLoggedIn) return <Login onLogin={() => setIsLoggedIn(true)} />
    return <AdminPanel />
}

export default App