import React, { useState, useEffect } from 'react';
import { User, Briefcase } from 'lucide-react';
import { api } from '../api';

interface AuthPageProps {
  onLogin: (role: 'user' | 'handyman', data?: any) => void;
}

export function AuthPage({ onLogin }: AuthPageProps) {
  const [isLogin, setIsLogin] = useState(true);
  const [role, setRole] = useState<'user' | 'handyman'>('user');
  
  // Form State
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  
  // Handyman Specific State
  const [category, setCategory] = useState('');
  const [hourlyRate, setHourlyRate] = useState('');
  const [description, setDescription] = useState('');

  // Categories from DB
  const [categories, setCategories] = useState<{ id: number; name: string }[]>([]);

  useEffect(() => {
    api.getCategories().then((data) => {
      if (Array.isArray(data)) {
        setCategories(data);
        if (data.length > 0 && !category) {
          setCategory(data[0].name);
        }
      }
    }).catch(console.error);
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    try {
      if (isLogin) {
        const response = await api.login({ email, password });
        if (response.token || response.access_token) {
          localStorage.setItem('token', response.token || response.access_token);
        }
        const backendUser = response.user || response;
        // Используем выбранную роль на UI (role) вместо бэкенд роли
        onLogin(role, backendUser);
      } else {
        // Registration
        // Backend ожидает числовое поле role (1 - клиент, 2 - мастер)
        const backendRole = role === 'handyman' ? 2 : 1;
        const userData = {
          name,
          email,
          password,
          role: backendRole,
          ...(role === 'handyman'
            ? {
                skill: category,
                hourlyRate: Number(hourlyRate),
                description: description,
              }
            : {}),
        };
        const response = await api.register(userData);
        if (response.token || response.access_token) {
          localStorage.setItem('token', response.token || response.access_token);
        }
        onLogin(role, response.user || response);
      }
    } catch (error) {
      console.error('Authentication failed:', error);
      const message = error instanceof Error ? error.message : 'Authentication failed. Please try again.';
      alert(message);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-100 p-4">
      <div className="w-full max-w-md bg-white border-2 border-gray-800 shadow-[8px_8px_0px_0px_rgba(31,41,55,1)] p-8">
        <div className="text-center mb-8">
          <div className="inline-block w-12 h-12 bg-gray-800 rounded-lg mb-4 flex items-center justify-center">
            <span className="text-white font-bold text-2xl">H</span>
          </div>
          <h1 className="text-3xl font-bold text-gray-900">{isLogin ? 'Добро пожаловать' : 'Создать аккаунт'}</h1>
          <p className="text-gray-500 mt-2">
            {isLogin ? 'Введите свои данные для доступа к аккаунту' : 'Зарегистрируйтесь, чтобы найти лучших мастеров'}
          </p>
        </div>

        {/* Role Selector - Visible on both Login/Register for demo purposes so we can easily switch roles */}
        <div className="flex bg-gray-100 p-1 rounded-lg mb-6">
          <button 
            type="button"
            onClick={() => setRole('user')}
            className={`flex-1 flex items-center justify-center py-2 rounded-md text-sm font-bold transition-all ${
              role === 'user' 
                ? 'bg-white text-gray-900 shadow-sm' 
                : 'text-gray-500 hover:text-gray-700'
            }`}
          >
            <User className="w-4 h-4 mr-2" />
            Пользователь
          </button>
          <button 
            type="button"
            onClick={() => setRole('handyman')}
            className={`flex-1 flex items-center justify-center py-2 rounded-md text-sm font-bold transition-all ${
              role === 'handyman' 
                ? 'bg-white text-gray-900 shadow-sm' 
                : 'text-gray-500 hover:text-gray-700'
            }`}
          >
            <Briefcase className="w-4 h-4 mr-2" />
            Мастер
          </button>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          {!isLogin && (
            <div>
              <label className="block text-sm font-bold text-gray-700 mb-1">Полное имя</label>
              <input 
                type="text" 
                value={name}
                onChange={(e) => setName(e.target.value)}
                className="w-full border-2 border-gray-300 p-2.5 rounded-md focus:border-gray-800 outline-none transition-colors"
                placeholder="Иван Иванов"
                required
              />
            </div>
          )}

          {/* Handyman Specific Fields */}
          {!isLogin && role === 'handyman' && (
            <>
              <div>
                <label className="block text-sm font-bold text-gray-700 mb-1">Специализация</label>
                <select 
                  value={category}
                  onChange={(e) => setCategory(e.target.value)}
                  className="w-full border-2 border-gray-300 p-2.5 rounded-md focus:border-gray-800 outline-none bg-white"
                >
                  {categories.length > 0 ? (
                    categories.map((cat) => (
                      <option key={cat.id} value={cat.name}>{cat.name}</option>
                    ))
                  ) : (
                    <option disabled>Загрузка...</option>
                  )}
                </select>
              </div>

              <div>
                <label className="block text-sm font-bold text-gray-700 mb-1">Почасовая ставка (руб.)</label>
                <input 
                  type="number" 
                  value={hourlyRate}
                  onChange={(e) => setHourlyRate(e.target.value)}
                  className="w-full border-2 border-gray-300 p-2.5 rounded-md focus:border-gray-800 outline-none transition-colors"
                  placeholder="50"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-bold text-gray-700 mb-1">Краткая биография</label>
                <textarea 
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                  className="w-full border-2 border-gray-300 p-2.5 rounded-md focus:border-gray-800 outline-none transition-colors resize-none"
                  placeholder="У меня 10 лет опыта..."
                  rows={2}
                  required
                />
              </div>
            </>
          )}

          <div>
            <label className="block text-sm font-bold text-gray-700 mb-1">Адрес электронной почты</label>
            <input 
              type="text" // Changed from 'email' to 'text' to allow non-email values like 'admin'
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full border-2 border-gray-300 p-2.5 rounded-md focus:border-gray-800 outline-none transition-colors"
              placeholder="you@example.com"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-bold text-gray-700 mb-1">Пароль</label>
            <input 
              type="password" 
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full border-2 border-gray-300 p-2.5 rounded-md focus:border-gray-800 outline-none transition-colors"
              placeholder="••••••••"
              required
            />
          </div>

          <button 
            type="submit" 
            className="w-full bg-gray-800 text-white font-bold py-3 px-4 rounded-md hover:bg-gray-700 active:transform active:translate-y-1 transition-all mt-4"
          >
            {isLogin ? 'Войти' : 'Присоединиться как ' + (role === 'handyman' ? 'Профи' : 'Пользователь')}
          </button>
        </form>

        <div className="mt-6 text-center text-sm text-gray-600">
          {isLogin ? "Нет аккаунта? " : "Уже есть аккаунт? "}
          <button 
            onClick={() => setIsLogin(!isLogin)} 
            className="font-bold text-gray-800 hover:underline"
          >
            {isLogin ? 'Зарегистрироваться' : 'Войти'}
          </button>
        </div>
      </div>
    </div>
  );
}
