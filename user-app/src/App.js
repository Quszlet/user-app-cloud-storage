import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import ProtectedRoute from "./Components/ProtectedRoute";

import LoginPage from "./Pages/LoginPage"
import RegistrationPage from "./Pages/RegistrationPage"
import HomePage from "./Pages/HomePage"

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/registration" element={<RegistrationPage />} />
        <Route
          path="/"
          element={<ProtectedRoute element={<HomePage />} redirectTo="/login" />}
        />
      </Routes>
    </Router>
  );
}

export default App;
