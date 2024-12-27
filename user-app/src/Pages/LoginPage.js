import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom'; // Импорт компонента Link
import '../CSS/Login.css'; // Подключение CSS для стилизации

const LoginPage = ({ element, redirectTo }) => {
  const navigate = useNavigate();
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [errors, setErrors] = useState({});
  const [authError, setAuthError] = useState('');

  const validateForm = () => {
    const newErrors = {};
    if (!username.trim()) {
      newErrors.username = 'Имя пользователя не должно быть пустым';
    }
    if (!password.trim()) {
      newErrors.password = 'Пароль не должен быть пустым';
    }
    return newErrors;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    // Сбрасываем предыдущие ошибки перед валидацией
    setErrors({});
    setAuthError('');

    const formErrors = validateForm();
    if (Object.keys(formErrors).length > 0) {
      setErrors(formErrors);
      return;
    }

    // Запрос к реальному API для проверки логина и пароля
    const isValid = await authenticateUser(username, password);

    if (!isValid) {
      setAuthError('Неверное имя пользователя или пароль');
    } else {
      setAuthError('');
      navigate("/")
    }
  };

  const authenticateUser = async (username, password) => {
    try {
        const obj = {login: username, password: password}
        const response = await fetch('http://localhost:8040/authentication', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(obj),
        });

        if (!response.ok) {
            console.log(response.body)
            return false;
        }

        const data = await response.json();
            // Сохранение токена или другой логики аутентификации
            localStorage.setItem('token', data.token);
            return true;
        } catch (error) {
            console.error('Ошибка при аутентификации:', error);
            return false;
        }
  };

  return (
    <div>
      <header className="header">
        <h1>Облачное хранилище</h1>
      </header>
      <div className="login-container">
        <form className="login-form" onSubmit={handleSubmit}>
          <h2>Вход</h2>

          {authError && <div className="auth-error">{authError}</div>}

          <div className="form-group">
            <label htmlFor="username">Имя пользователя</label>
            <input
              type="text"
              id="username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className={errors.username ? 'input-error' : ''}
            />
            {errors.username && <small className="error-text">{errors.username}</small>}
          </div>

          <div className="form-group">
            <label htmlFor="password">Пароль</label>
            <input
              type="password"
              id="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className={errors.password ? 'input-error' : ''}
            />
            {errors.password && <small className="error-text">{errors.password}</small>}
          </div>

          <button type="submit">Войти</button>
        </form>

        <p className="register-link">
            Нет аккаунта? <Link to="/registration">Зарегистрируйтесь</Link>
        </p>

      </div>
    </div>
  );
};

export default LoginPage;
