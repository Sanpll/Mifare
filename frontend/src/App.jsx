import React, { useState, useEffect } from 'react'

const decodeJWT = (token) => {
    try {
        const base64Url = token.split('.')[1]
        const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/')
        const jsonPayload = decodeURIComponent(
            atob(base64)
                .split('')
                .map(c => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
                .join('')
        )
        return JSON.parse(jsonPayload)
    } catch {
        return null
    }
}

let authToken = localStorage.getItem('token') || null
let currentUserId = null
let currentUsername = null
let isAdmin = false

if (authToken) {
    const payload = decodeJWT(authToken)
    if (payload) {
        currentUserId = payload.user_id
        currentUsername = payload.username
        isAdmin = payload.is_admin === true
    }
}

// API-клиент
const request = async (endpoint, options = {}) => {
    const headers = { 'Content-Type': 'application/json', ...options.headers }
    if (authToken) headers['Authorization'] = `Bearer ${authToken}`
    const res = await fetch(endpoint, { ...options, headers })
    if (res.status === 401 && authToken) {
        localStorage.removeItem('token')
        authToken = null
        window.location.reload()
    }
    return res
}

// Логин
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
            const payload = decodeJWT(data.token)
            if (payload) {
                authToken = data.token
                currentUserId = payload.user_id
                currentUsername = payload.username
                isAdmin = payload.is_admin === true
                localStorage.setItem('token', authToken)
                onLogin()
            } else {
                setError('Ошибка токена')
            }
        } else {
            setError('Неверное имя или пароль')
        }
    }

    return (
        <div className="min-h-screen flex items-center justify-center bg-gray-50">
            <div className="bg-white p-8 rounded shadow-md w-96">
                <h1 className="text-2xl font-bold mb-6">Вход</h1>
                <form onSubmit={handleSubmit}>
                    <input type="text" placeholder="Имя" value={username} onChange={e => setUsername(e.target.value)} className="border p-2 w-full mb-4 rounded" required />
                    <input type="password" placeholder="Пароль" value={password} onChange={e => setPassword(e.target.value)} className="border p-2 w-full mb-4 rounded" required />
                    {error && <p className="text-red-500 mb-4">{error}</p>}
                    <button type="submit" className="bg-blue-600 text-white py-2 px-4 rounded w-full">Войти</button>
                </form>
            </div>
        </div>
    )
}

// Общая заготовка для всех таблиц
const Table = ({ items, fields, getId = (item) => item.id, actions = null }) => (
    <table className="w-full border-collapse border">
        <thead><tr className="bg-gray-100">{fields.map(f => <th key={f.name} className="border p-2">{f.label}</th>)}<th className="border p-2">Действия</th></tr></thead>
        <tbody>
            {items.map(item => (
                <tr key={getId(item)}>
                    {fields.map(f => <td key={f.name} className="border p-2">{item[f.name]?.toString()}</td>)}
                    <td className="border p-2">{actions ? actions(item) : ''}</td>
                </tr>
            ))}
        </tbody>
    </table>
)

// Профиль пользователя
const Profile = () => {
    const [user, setUser] = useState(null)
    useEffect(() => {
        request(`/api/v1/users/${currentUserId}`).then(res => res.ok && res.json().then(setUser))
    }, [])
    if (!user) return <div className="p-4">Загрузка...</div>
    return (
        <div className="p-4">
            <h2 className="text-xl font-bold mb-4">Мой профиль</h2>
            <p><strong>ID:</strong> {user.id}</p>
            <p><strong>Имя:</strong> {user.username}</p>
            <p><strong>Администратор:</strong> {user.is_admin ? 'Да' : 'Нет'}</p>
            <p><strong>Дата регистрации:</strong> {new Date(user.created_at).toLocaleString()}</p>
        </div>
    )
}

// ---- Карты пользователя (только его карты, создание, без редактирования/удаления) ----
const UserCards = () => {
    const [cards, setCards] = useState([])
    const [form, setForm] = useState({ number: '', balance: '', is_blocked: false, key_value: '' })

    const fetchCards = async () => {
        const res = await request('/api/v1/cards')
        if (res.ok) {
            const data = await res.json()
            // оставляем только карты текущего пользователя
            setCards((data.cards || []).filter(c => c.owner_name === currentUsername))
        }
    }

    useEffect(() => { fetchCards() }, [])

    const handleCreate = async (e) => {
        e.preventDefault()
        const payload = { ...form, balance: parseFloat(form.balance) }
        const res = await request('/api/v1/cards', { method: 'POST', body: JSON.stringify(payload) })
        if (res.ok) {
            fetchCards()
            setForm({ number: '', balance: '', is_blocked: false, key_value: '' })
        } else {
            alert('Ошибка создания')
        }
    }

    const fields = [{ name: 'number', label: 'Номер' }, { name: 'balance', label: 'Баланс' }, { name: 'is_blocked', label: 'Заблокирована' }, { name: 'owner_name', label: 'Владелец' }, { name: 'key_value', label: 'Ключ' }]
    return (
        <div className="p-4">
            <h2 className="text-xl font-bold mb-4">Мои карты</h2>
            <form onSubmit={handleCreate} className="mb-6 flex flex-wrap gap-2">
                <input type="text" placeholder="Номер (hex)" value={form.number} onChange={e => setForm({ ...form, number: e.target.value })} className="border p-2 rounded" required />
                <input type="number" step="0.01" placeholder="Баланс" value={form.balance} onChange={e => setForm({ ...form, balance: e.target.value })} className="border p-2 rounded" required />
                <input type="text" placeholder="Ключ (12 hex)" value={form.key_value} onChange={e => setForm({ ...form, key_value: e.target.value })} className="border p-2 rounded" required />
                <button type="submit" className="bg-green-600 text-white px-4 py-2 rounded">Создать</button>
            </form>
            {cards.length === 0 ? <p>Нет карт</p> : <Table items={cards} fields={fields} actions={() => ''} />}
        </div>
    )
}

// Карты для администратора
const AdminCards = () => {
    const [cards, setCards] = useState([])
    const [editingId, setEditingId] = useState(null)
    const [form, setForm] = useState({})

    const fetchCards = async () => {
        const res = await request('/api/v1/cards')
        if (res.ok) setCards((await res.json()).cards || [])
    }
    useEffect(() => { fetchCards() }, [])

    const handleDelete = async (id) => {
        if (confirm('Удалить карту?')) {
            const res = await request(`/api/v1/cards/${id}`, { method: 'DELETE' })
            if (res.ok) fetchCards()
        }
    }

    const handleUpdate = async (e) => {
        e.preventDefault()
        const res = await request(`/api/v1/cards/${editingId}`, { method: 'PUT', body: JSON.stringify(form) })
        if (res.ok) {
            fetchCards()
            setEditingId(null)
            setForm({})
        } else {
            alert('Ошибка обновления')
        }
    }

    const startEdit = (card) => {
        setEditingId(card.id)
        setForm({ balance: card.balance, is_blocked: card.is_blocked, owner_name: card.owner_name, key_value: card.key_value })
    }

    const fields = [{ name: 'number', label: 'Номер' }, { name: 'balance', label: 'Баланс' }, { name: 'is_blocked', label: 'Заблокирована' }, { name: 'owner_name', label: 'Владелец' }, { name: 'key_value', label: 'Ключ' }]
    return (
        <div className="p-4">
            <h2 className="text-xl font-bold mb-4">Все карты</h2>
            {editingId && (
                <form onSubmit={handleUpdate} className="mb-6 flex flex-wrap gap-2">
                    <input type="number" step="0.01" placeholder="Баланс" value={form.balance || ''} onChange={e => setForm({ ...form, balance: parseFloat(e.target.value) })} className="border p-2 rounded" />
                    <label className="flex items-center gap-1"><input type="checkbox" checked={form.is_blocked || false} onChange={e => setForm({ ...form, is_blocked: e.target.checked })} /> Заблокирована</label>
                    <input type="text" placeholder="Владелец" value={form.owner_name || ''} onChange={e => setForm({ ...form, owner_name: e.target.value })} className="border p-2 rounded" />
                    <input type="text" placeholder="Ключ (12 hex)" value={form.key_value || ''} onChange={e => setForm({ ...form, key_value: e.target.value })} className="border p-2 rounded" />
                    <button type="submit" className="bg-blue-600 text-white px-4 py-2 rounded">Обновить</button>
                    <button type="button" onClick={() => { setEditingId(null); setForm({}) }} className="bg-gray-400 text-white px-4 py-2 rounded">Отмена</button>
                </form>
            )}
            <Table items={cards} fields={fields} actions={(card) => (
                <>
                    <button onClick={() => startEdit(card)} className="bg-yellow-500 text-white px-2 py-1 rounded mr-2">Ред.</button>
                    <button onClick={() => handleDelete(card.id)} className="bg-red-500 text-white px-2 py-1 rounded">Уд.</button>
                </>
            )} />
        </div>
    )
}

// Проверка транзакции
const TransactionAuth = () => {
    const [cardNumber, setCardNumber] = useState('')
    const [price, setPrice] = useState('')
    const [terminalSN, setTerminalSN] = useState('')
    const [result, setResult] = useState(null)
    const [userCards, setUserCards] = useState([])

    useEffect(() => {
        if (!isAdmin) {
            request('/api/v1/cards').then(res => res.ok && res.json().then(data => {
                setUserCards((data.cards || []).filter(c => c.owner_name === currentUsername))
            }))
        }
    }, [])

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
                {!isAdmin ? (
                    <select value={cardNumber} onChange={e => setCardNumber(e.target.value)} className="border p-2 rounded">
                        <option value="">Выберите карту</option>
                        {userCards.map(c => <option key={c.id} value={c.number}>{c.number}</option>)}
                    </select>
                ) : (
                    <input placeholder="Номер карты" value={cardNumber} onChange={e => setCardNumber(e.target.value)} className="border p-2 rounded" />
                )}
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

// Админ панель
const AdminPanel = () => {
    const [activeTab, setActiveTab] = useState(isAdmin ? 'users' : 'profile')
    const logout = () => {
        localStorage.removeItem('token')
        authToken = null
        window.location.reload()
    }

    const adminTabs = [
        { id: 'users', label: 'Пользователи', component: () => <Table items={[]} fields={[{ name: 'username', label: 'Имя' }]} actions={() => ''} /> }
    ]

    const UsersTable = () => {
        const [users, setUsers] = useState([])
        useEffect(() => {
            request('/api/v1/users').then(res => res.ok && res.json().then(data => setUsers(data.users || [])))
        }, [])
        return (
            <div className="p-4">
                <h2 className="text-xl font-bold mb-4">Пользователи</h2>
                <Table items={users} fields={[{ name: 'username', label: 'Имя' }]} actions={() => ''} />
            </div>
        )
    }

    const adminTabsFinal = [
        { id: 'users', label: 'Пользователи', component: UsersTable },
        {
            id: 'keys', label: 'Ключи', component: () => {
                const [keys, setKeys] = useState([])
                useEffect(() => { request('/api/v1/keys').then(res => res.ok && res.json().then(data => setKeys(data.keys || []))) }, [])
                return <div className="p-4"><h2 className="text-xl font-bold mb-4">Ключи</h2><Table items={keys} fields={[{ name: 'value', label: 'Значение' }, { name: 'type', label: 'Тип' }, { name: 'description', label: 'Описание' }]} actions={() => ''} /></div>
            }
        },
        { id: 'cards', label: 'Карты', component: AdminCards },
        {
            id: 'terminals', label: 'Терминалы', component: () => {
                const [terminals, setTerminals] = useState([])
                useEffect(() => { request('/api/v1/terminals').then(res => res.ok && res.json().then(data => setTerminals(data.terminals || []))) }, [])
                return <div className="p-4"><h2 className="text-xl font-bold mb-4">Терминалы</h2><Table items={terminals} fields={[{ name: 'serial_number', label: 'Серийный номер' }, { name: 'address', label: 'Адрес' }, { name: 'name', label: 'Название' }]} actions={() => ''} /></div>
            }
        },
        {
            id: 'transactions', label: 'Транзакции', component: () => {
                const [transactions, setTransactions] = useState([])
                useEffect(() => { request('/api/v1/transactions').then(res => res.ok && res.json().then(data => setTransactions(data.transactions || []))) }, [])
                return <div className="p-4"><h2 className="text-xl font-bold mb-4">Транзакции</h2><Table items={transactions} fields={[{ name: 'card_number', label: 'Номер карты' }, { name: 'price', label: 'Сумма' }, { name: 'terminal_serial_number', label: 'Терминал' }]} actions={() => ''} /></div>
            }
        },
        { id: 'auth', label: 'Проверка транзакции', component: TransactionAuth }
    ]

    const userTabs = [
        { id: 'profile', label: 'Профиль', component: Profile },
        { id: 'cards', label: 'Мои карты', component: UserCards },
        { id: 'auth', label: 'Проверка транзакции', component: TransactionAuth }
    ]

    const tabs = isAdmin ? adminTabsFinal : userTabs
    const ActiveComponent = tabs.find(t => t.id === activeTab)?.component || tabs[0].component

    return (
        <div>
            <div className="bg-gray-800 text-white p-4 flex justify-between">
                <div className="flex gap-4 flex-wrap">
                    {tabs.map(tab => (
                        <button key={tab.id} onClick={() => setActiveTab(tab.id)} className={`px-3 py-1 rounded ${activeTab === tab.id ? 'bg-blue-600' : 'hover:bg-gray-700'}`}>
                            {tab.label}
                        </button>
                    ))}
                </div>
                <button onClick={logout} className="bg-red-600 px-3 py-1 rounded">Выйти</button>
            </div>
            <div className="p-4"><ActiveComponent /></div>
        </div>
    )
}

const App = () => {
    const [isLoggedIn, setIsLoggedIn] = useState(!!authToken)
    if (!isLoggedIn) return <Login onLogin={() => setIsLoggedIn(true)} />
    return <AdminPanel />
}

export default App