import React from "react";
import { Link, useNavigate } from "react-router-dom";

export default function Header() {
  const token = localStorage.getItem("token");
  const navigate = useNavigate();

  const handleLogout = () => {
    localStorage.removeItem("token");
    navigate("/login");
  };

  return (
    <header style={styles.header}>
      <div style={styles.container}>
        <Link to={token ? "/home" : "/"} style={styles.brand}>Challenge Excel</Link>
        <nav>
          {token ? (
            <>
              <Link to="/home" style={styles.link}>Home</Link>
              <Link to="/import" style={styles.link}>Importar</Link>
              <button onClick={handleLogout} style={styles.btn}>Logout</button>
            </>
          ) : (
            <>
              <Link to="/login" style={styles.link}>Login</Link>
              <Link to="/register" style={styles.link}>Register</Link>
            </>
          )}
        </nav>
      </div>
    </header>
  );
}

const styles = {
  header: { background: "#ffffff", borderBottom: "1px solid #eee", boxShadow: "0 1px 0 rgba(0,0,0,0.03)" },
  container: { maxWidth: 1100, margin: "0 auto", padding: "0.6rem 1rem", display: "flex", alignItems: "center", justifyContent: "space-between" },
  brand: { fontWeight: 700, color: "#333", textDecoration: "none", fontSize: "1.05rem" },
  link: { marginLeft: 12, color: "#4a90e2", textDecoration: "none", fontWeight: 600 },
  btn: { marginLeft: 12, padding: "6px 10px", borderRadius: 8, border: "none", background: "#ef4444", color: "#fff", cursor: "pointer", fontWeight: 600 },
};
