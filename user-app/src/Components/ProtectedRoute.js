import React, { useEffect, useState } from 'react';
import { Navigate } from 'react-router-dom';

const ProtectedRoute = ({ element, redirectTo }) => {
  const [isAuthenticated, setIsAuthenticated] = useState(null); // null для состояния загрузки

  const validateToken = async (token) => {
    try {
      const response = await fetch('http://localhost:8040/validation', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
      });
  
      if (!response.ok) {
        throw new Error('Invalid token');
      }
  
      return true; // Токен валиден
    } catch (error) {
      console.error(error);
      return false; // Токен не валиден
    }
  };

  useEffect(() => {
    const checkAuth = async () => {
      const token = localStorage.getItem('token'); // Получение токена из localStorage или другого источника
      if (token) {
        const isValid = await validateToken(token); // Валидация токена
        setIsAuthenticated(isValid);
      } else {
        setIsAuthenticated(false);
      }
    };

    checkAuth();
  }, []);

  if (isAuthenticated === null) {
    return <div>Проверка авторизации...</div>; // Спиннер или индикатор загрузки
  }

  return isAuthenticated ? element : <Navigate to={redirectTo} />;
};

export default ProtectedRoute;
