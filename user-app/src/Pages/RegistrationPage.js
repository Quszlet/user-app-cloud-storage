import React, { useState } from 'react';
import { Link, useNavigate  } from 'react-router-dom'; 
import '../CSS/Register.css'; // Подключение CSS для стилизации

const RegistrationPage = ({ element, redirectTo }) => {
  const navigate = useNavigate();
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [errors, setErrors] = useState({});
  const [registerError, setRegisterError] = useState('');

  const validateForm = () => {
    const newErrors = {};

    if (!username.trim()) {
      newErrors.username = 'Имя пользователя не должно быть пустым';
    }
    if (!email.trim()) {
      newErrors.email = 'Email не должен быть пустым';
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
      newErrors.email = 'Неверный формат email';
    }
    if (!password.trim()) {
      newErrors.password = 'Пароль не должен быть пустым';
    }
    if (password !== confirmPassword) {
      newErrors.confirmPassword = 'Пароли не совпадают';
    }

    return newErrors;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    // Сбрасываем предыдущие ошибки перед валидацией
    setErrors({});
    setRegisterError('');

    const formErrors = validateForm();
    if (Object.keys(formErrors).length > 0) {
      setErrors(formErrors);
      return;
    }

    // Запрос к реальному API для регистрации пользователя
    const isRegistered = await registerUser(username, email, password);

    if (!isRegistered) {
      setRegisterError('Ошибка при регистрации');
    } else {
      setRegisterError('');
      alert('Успешная регистрация!'); // Замените на перенаправление
    }
  };

  const registerUser = async (username, email, password) => {
    try {
      const user = {
        login: username,
        password: password,
        email: email
      }
      console.log(user)
      const response = await fetch('http://localhost:8040/users/create', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(user),
      });

      if (!response.ok) {
        return false;
      }

      const data = await response.json();
      // Дополнительная логика после успешной регистрации
      const bodyDir = {
        size: 0,
        name: "home",
        directory_id: null,
        user_id: data.id
      }

      console.log(bodyDir)

      const responseCreate = await fetch('http://localhost:8040/directories/create', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(bodyDir),
      });

      if (!responseCreate.ok) {
        return false;
      }

      navigate("/login")
      return true;
    } catch (error) {
      console.error('Ошибка при регистрации:', error);
      return false;
    }
  };

  return (
    <div>
      <header className="header">
        <h1>Облачное хранилище</h1>
      </header>
      <div className="register-container">
        <form className="register-form" onSubmit={handleSubmit}>
          <h2>Регистрация</h2>

          {registerError && <div className="register-error">{registerError}</div>}

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
            <label htmlFor="email">Email</label>
            <input
              type="email"
              id="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className={errors.email ? 'input-error' : ''}
            />
            {errors.email && <small className="error-text">{errors.email}</small>}
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

          <div className="form-group">
            <label htmlFor="confirmPassword">Подтверждение пароля</label>
            <input
              type="password"
              id="confirmPassword"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              className={errors.confirmPassword ? 'input-error' : ''}
            />
            {errors.confirmPassword && <small className="error-text">{errors.confirmPassword}</small>}
          </div>

          <button type="submit">Зарегистрироваться</button>
        </form>

        <p className="register-link">
            Есть аккаунта? <Link to="/login">Выполните вход</Link>
        </p>
      </div>
    </div>
  );
};

export default RegistrationPage;
