import React, { useState, useEffect } from 'react';

const decodeJWT = (token) => {
  try {
    const base64Url = token.split('.')[1];
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const jsonPayload = decodeURIComponent(
      atob(base64)
        .split('')
        .map((c) => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
        .join('')
    );
    return JSON.parse(jsonPayload);
  } catch {
    return null;
  }
};

let authToken = localStorage.getItem('token') || null;
let currentUserId = null;
let currentUsername = null;
let isAdmin = false;

if (authToken) {
  const payload = decodeJWT(authToken);
  if (payload) {
    currentUserId = payload.user_id;
    currentUsername = payload.username;
    isAdmin = payload.is_admin === true;
  }
}

const request = async (endpoint, options = {}) => {
  const headers = { 'Content-Type': 'application/json', ...options.headers };
  if (authToken) headers['Authorization'] = `Bearer ${authToken}`;
  const res = await fetch(endpoint, { ...options, headers });
  if (res.status === 401 && authToken) {
    localStorage.removeItem('token');
    authToken = null;
    window.location.reload();
  }
  return res;
};

// Общие компоненты
const Login = ({ onLogin }) => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    const res = await fetch('/auth/sign-in', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, password }),
    });
    if (res.ok) {
      const data = await res.json();
      const payload = decodeJWT(data.token);
      if (payload) {
        authToken = data.token;
        currentUserId = payload.user_id;
        currentUsername = payload.username;
        isAdmin = payload.is_admin === true;
        localStorage.setItem('token', authToken);
        onLogin();
      } else {
        setError('Ошибка токена');
      }
    } else {
      setError('Неверное имя или пароль');
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="bg-white p-8 rounded shadow-md w-96">
        <h1 className="text-2xl font-bold mb-6">Вход</h1>
        <form onSubmit={handleSubmit}>
          <input
            type="text"
            placeholder="Имя"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            className="border p-2 w-full mb-4 rounded"
            required
          />
          <input
            type="password"
            placeholder="Пароль"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="border p-2 w-full mb-4 rounded"
            required
          />
          {error && <p className="text-red-500 mb-4">{error}</p>}
          <button type="submit" className="bg-blue-600 text-white py-2 px-4 rounded w-full">
            Войти
          </button>
        </form>
      </div>
    </div>
  );
};

const Table = ({ items, fields, getId = (item) => item.id, actions = null }) => {
  const hasActions = actions !== null && actions !== undefined;
  return (
    <table className="w-full border-collapse border">
      <thead>
        <tr className="bg-gray-100">
          {fields.map((f) => (
            <th key={f.name} className="border p-2">
              {f.label}
            </th>
          ))}
          {hasActions && <th className="border p-2">Действия</th>}
        </tr>
      </thead>
      <tbody>
        {items.map((item) => (
          <tr key={getId(item)}>
            {fields.map((f) => (
              <td key={f.name} className="border p-2">
                {item[f.name]?.toString()}
              </td>
            ))}
            {hasActions && <td className="border p-2">{actions(item)}</td>}
          </tr>
        ))}
      </tbody>
    </table>
  );
};

const Profile = () => {
  const [user, setUser] = useState(null);
  useEffect(() => {
    request(`/api/v1/users/${currentUserId}`).then((res) =>
      res.ok && res.json().then(setUser)
    );
  }, []);
  if (!user) return <div className="p-4">Загрузка...</div>;
  return (
    <div className="p-4">
      <h2 className="text-xl font-bold mb-4">Мой профиль</h2>
      <p><strong>ID:</strong> {user.id}</p>
      <p><strong>Имя:</strong> {user.username}</p>
      <p><strong>Администратор:</strong> {user.is_admin ? 'Да' : 'Нет'}</p>
      <p><strong>Дата регистрации:</strong> {new Date(user.created_at).toLocaleString()}</p>
    </div>
  );
};

const UserCards = () => {
  const [cards, setCards] = useState([]);
  useEffect(() => {
    request('/api/v1/cards').then((res) =>
      res.ok && res.json().then((data) =>
        setCards((data.cards || []).filter((c) => c.owner_name === currentUsername))
      )
    );
  }, []);
  const fields = [
    { name: 'number', label: 'Номер' },
    { name: 'balance', label: 'Баланс' },
    { name: 'is_blocked', label: 'Заблокирована' },
    { name: 'owner_name', label: 'Владелец' },
    { name: 'key_value', label: 'Ключ' },
  ];
  return (
    <div className="p-4">
      <h2 className="text-xl font-bold mb-4">Мои карты</h2>
      {cards.length === 0 ? <p>Нет карт</p> : <Table items={cards} fields={fields} />}
    </div>
  );
};

const TransactionAuth = () => {
  const [cardNumber, setCardNumber] = useState('');
  const [price, setPrice] = useState('');
  const [terminalSN, setTerminalSN] = useState('');
  const [result, setResult] = useState(null);
  const [userCards, setUserCards] = useState([]);

  useEffect(() => {
    if (!isAdmin) {
      request('/api/v1/cards').then((res) =>
        res.ok && res.json().then((data) =>
          setUserCards((data.cards || []).filter((c) => c.owner_name === currentUsername))
        )
      );
    }
  }, []);

  const checkAuth = async () => {
    const res = await request('/api/v1/terminal/auth', {
      method: 'POST',
      body: JSON.stringify({ card_number: cardNumber, price, terminal_serial_number: terminalSN }),
    });
    const data = await res.json();
    setResult(data);
  };

  return (
    <div className="p-4">
      <h2 className="text-xl font-bold mb-4">Авторизация транзакции</h2>
      <div className="mb-4 flex gap-2 flex-wrap">
        {!isAdmin ? (
          <select value={cardNumber} onChange={(e) => setCardNumber(e.target.value)} className="border p-2 rounded">
            <option value="">Выберите карту</option>
            {userCards.map((c) => (
              <option key={c.id} value={c.number}>{c.number}</option>
            ))}
          </select>
        ) : (
          <input
            placeholder="Номер карты"
            value={cardNumber}
            onChange={(e) => setCardNumber(e.target.value)}
            className="border p-2 rounded"
          />
        )}
        <input
          placeholder="Сумма"
          value={price}
          onChange={(e) => setPrice(e.target.value)}
          className="border p-2 rounded"
        />
        <input
          placeholder="Серийный номер терминала"
          value={terminalSN}
          onChange={(e) => setTerminalSN(e.target.value)}
          className="border p-2 rounded"
        />
        <button onClick={checkAuth} className="bg-blue-600 text-white px-4 py-2 rounded">
          Проверить
        </button>
      </div>
      {result && (
        <div className={`p-3 rounded ${result.authorized ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}`}>
          {result.message}
        </div>
      )}
    </div>
  );
};

// Универсальный CRUD-менеджер
const CrudManager = ({
  endpoint,
  fields,           // [{ name, label, type?: 'text'|'checkbox', placeholder?, required?, disabledOnEdit? }]
  tableFields,      // [{ name, label }]
  title,
  getId = (item) => item.id,
  excludeOnUpdate = [],   // поля, не отправляемые при обновлении
  allowCreate = true,     // разрешить создание новых записей
  readonly = false,       // если true – только таблица без действий и формы
}) => {
  const [items, setItems] = useState([]);
  const [editingId, setEditingId] = useState(null);

  const initialForm = fields.reduce((acc, f) => ({
    ...acc,
    [f.name]: f.type === 'checkbox' ? false : '',
  }), {});
  const [form, setForm] = useState(initialForm);

  const resource = endpoint.replace('/api/v1/', '');
  const fetchItems = async () => {
    const res = await request(endpoint);
    if (res.ok) {
      const data = await res.json();
      setItems(data[resource] || []);
    }
  };

  useEffect(() => { fetchItems(); }, [endpoint]);

  const handleDelete = async (id) => {
    if (confirm(`Удалить ${title.slice(0, -1)}?`)) {
      const res = await request(`${endpoint}/${id}`, { method: 'DELETE' });
      if (res.ok) fetchItems();
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    const url = editingId ? `${endpoint}/${editingId}` : endpoint;
    const method = editingId ? 'PUT' : 'POST';
    let payload = { ...form };
    if (editingId) {
      excludeOnUpdate.forEach((name) => delete payload[name]);
    }
    const res = await request(url, { method, body: JSON.stringify(payload) });
    if (res.ok) {
      fetchItems();
      setEditingId(null);
      setForm(initialForm);
    } else {
      alert('Ошибка');
    }
  };

  const startEdit = (item) => {
    setEditingId(item.id);
    const newForm = {};
    fields.forEach((f) => {
      newForm[f.name] = item[f.name] !== undefined ? item[f.name] : (f.type === 'checkbox' ? false : '');
    });
    setForm(newForm);
  };

  const cancelEdit = () => {
    setEditingId(null);
    setForm(initialForm);
  };

  const showForm = !readonly && allowCreate;
  const showActions = !readonly;

  return (
    <div className="p-4">
      <h2 className="text-xl font-bold mb-4">{title}</h2>
      {showForm && (
        <form onSubmit={handleSubmit} className="mb-6 flex flex-wrap gap-2">
          {fields.map((f) => {
            if (f.type === 'checkbox') {
              return (
                <label key={f.name} className="flex items-center gap-1">
                  <input
                    type="checkbox"
                    checked={form[f.name] || false}
                    onChange={(e) => setForm({ ...form, [f.name]: e.target.checked })}
                    disabled={!!editingId && f.disabledOnEdit}
                  />
                  {f.label}
                </label>
              );
            }
            return (
              <input
                key={f.name}
                type="text"
                placeholder={f.placeholder || f.label}
                value={form[f.name] || ''}
                onChange={(e) => setForm({ ...form, [f.name]: e.target.value })}
                className="border p-2 rounded"
                required={f.required && !editingId}
                disabled={!!editingId && f.disabledOnEdit}
              />
            );
          })}
          <button type="submit" className="bg-green-600 text-white px-4 py-2 rounded">
            {editingId ? 'Обновить' : 'Создать'}
          </button>
          {editingId && (
            <button type="button" onClick={cancelEdit} className="bg-gray-400 text-white px-4 py-2 rounded">
              Отмена
            </button>
          )}
        </form>
      )}
      <Table
        items={items}
        fields={tableFields}
        actions={
          showActions
            ? (item) => (
                <>
                  <button onClick={() => startEdit(item)} className="bg-yellow-500 text-white px-2 py-1 rounded mr-2">
                    Ред.
                  </button>
                  <button onClick={() => handleDelete(item.id)} className="bg-red-500 text-white px-2 py-1 rounded">
                    Уд.
                  </button>
                </>
              )
            : null
        }
        getId={getId}
      />
    </div>
  );
};

// Разделы для админа
const AdminCards = () => (
  <CrudManager
    endpoint="/api/v1/cards"
    title="Все карты"
    fields={[
      { name: 'number', label: 'Номер', disabledOnEdit: true, required: true },
      { name: 'balance', label: 'Баланс', required: true },
      { name: 'is_blocked', label: 'Заблокирована', type: 'checkbox' },
      { name: 'owner_name', label: 'Владелец', required: true },
      { name: 'key_value', label: 'Ключ (12 hex)', required: true },
    ]}
    tableFields={[
      { name: 'number', label: 'Номер' },
      { name: 'balance', label: 'Баланс' },
      { name: 'is_blocked', label: 'Заблокирована' },
      { name: 'owner_name', label: 'Владелец' },
      { name: 'key_value', label: 'Ключ' },
    ]}
    excludeOnUpdate={['number']}
  />
);

const AdminUsers = () => (
  <CrudManager
    endpoint="/api/v1/users"
    title="Пользователи"
    fields={[{ name: 'username', label: 'Имя', required: true }]}
    tableFields={[{ name: 'username', label: 'Имя' }]}
    allowCreate={false} // создание пользователей не предусмотрено
  />
);

const AdminKeys = () => (
  <CrudManager
    endpoint="/api/v1/keys"
    title="Ключи"
    fields={[
      { name: 'value', label: 'Значение (12 hex)', required: true },
      { name: 'type', label: 'Тип (A/B)', required: true },
      { name: 'description', label: 'Описание', required: true },
    ]}
    tableFields={[
      { name: 'value', label: 'Значение' },
      { name: 'type', label: 'Тип' },
      { name: 'description', label: 'Описание' },
    ]}
  />
);

const TerminalsList = ({ admin = false }) => (
  <CrudManager
    endpoint="/api/v1/terminals"
    title={admin ? 'Все терминалы' : 'Терминалы'}
    fields={[
      { name: 'serial_number', label: 'Серийный номер', required: true },
      { name: 'address', label: 'Адрес', required: true },
      { name: 'name', label: 'Название', required: true },
    ]}
    tableFields={[
      { name: 'serial_number', label: 'Серийный номер' },
      { name: 'address', label: 'Адрес' },
      { name: 'name', label: 'Название' },
    ]}
    readonly={!admin}
    allowCreate={admin}
  />
);

const TransactionsList = () => {
  const [transactions, setTransactions] = useState([]);
  useEffect(() => {
    request('/api/v1/transactions').then((res) =>
      res.ok && res.json().then((data) => setTransactions(data.transactions || []))
    );
  }, []);
  const fields = [
    { name: 'card_number', label: 'Номер карты' },
    { name: 'price', label: 'Сумма' },
    { name: 'terminal_serial_number', label: 'Терминал' },
  ];
  return (
    <div className="p-4">
      <h2 className="text-xl font-bold mb-4">Транзакции</h2>
      <Table items={transactions} fields={fields} />
    </div>
  );
};

// Админ панель с вкладками
const AdminPanel = () => {
  const [activeTab, setActiveTab] = useState(isAdmin ? 'users' : 'profile');

  const logout = () => {
    localStorage.removeItem('token');
    authToken = null;
    window.location.reload();
  };

  const adminTabs = [
    { id: 'users', label: 'Пользователи', component: AdminUsers },
    { id: 'keys', label: 'Ключи', component: AdminKeys },
    { id: 'cards', label: 'Карты', component: AdminCards },
    { id: 'terminals', label: 'Терминалы', component: () => <TerminalsList admin /> },
    { id: 'transactions', label: 'Транзакции', component: TransactionsList },
    { id: 'auth', label: 'Проверка транзакции', component: TransactionAuth },
  ];

  const userTabs = [
    { id: 'profile', label: 'Профиль', component: Profile },
    { id: 'cards', label: 'Мои карты', component: UserCards },
    { id: 'terminals', label: 'Терминалы', component: () => <TerminalsList admin={false} /> },
    { id: 'auth', label: 'Проверка транзакции', component: TransactionAuth },
  ];

  const tabs = isAdmin ? adminTabs : userTabs;
  const ActiveComponent = tabs.find((t) => t.id === activeTab)?.component || tabs[0].component;

  return (
    <div>
      <div className="bg-gray-800 text-white p-4 flex justify-between">
        <div className="flex gap-4 flex-wrap">
          {tabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={`px-3 py-1 rounded ${activeTab === tab.id ? 'bg-blue-600' : 'hover:bg-gray-700'}`}
            >
              {tab.label}
            </button>
          ))}
        </div>
        <button onClick={logout} className="bg-red-600 px-3 py-1 rounded">
          Выйти
        </button>
      </div>
      <div className="p-4">
        <ActiveComponent />
      </div>
    </div>
  );
};

// Корневой компонент
const App = () => {
  const [isLoggedIn, setIsLoggedIn] = useState(!!authToken);
  if (!isLoggedIn) return <Login onLogin={() => setIsLoggedIn(true)} />;
  return <AdminPanel />;
};

export default App;