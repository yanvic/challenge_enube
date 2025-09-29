// src/App.js
import React from "react";
import { Routes, Route, Navigate, useLocation } from "react-router-dom";
import Login from "./pages/Login";
import Register from "./pages/Register";
import Home from "./pages/Home";
import ImportPage from "./pages/Import";
import Header from "./components/Header";
import ProtectedRoute from "./components/ProtectedRoute";

function App() {
  const token = localStorage.getItem("token");
  const location = useLocation();

  if (location.pathname === "/") {
    if (!token) return <Navigate to="/register" replace />;
  }

  return (
    <div style={{ fontFamily: "Inter, Roboto, Arial, sans-serif" }}>
      <Header />
      <Routes>
        <Route path="/" element={<Navigate to={token ? "/home" : "/register"} replace />} />
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />

        <Route
          path="/home"
          element={
            <ProtectedRoute>
              <Home />
            </ProtectedRoute>
          }
        />

        <Route
          path="/import"
          element={
            <ProtectedRoute>
              <ImportPage />
            </ProtectedRoute>
          }
        />

        {/* fallback */}
        <Route path="*" element={<Navigate to={token ? "/home" : "/register"} replace />} />
      </Routes>
    </div>
  );
}

export default App;
