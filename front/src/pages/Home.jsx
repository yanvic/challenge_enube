import React, { useEffect, useState } from "react";
import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer, PieChart, Pie, Cell } from "recharts";

export default function Home() {
  const [metrics, setMetrics] = useState(null);
  const [products, setProducts] = useState([]);
  const [customers, setCustomers] = useState([]);
  const [filterMonth, setFilterMonth] = useState("");

  const COLORS = ["#2563eb", "#16a34a", "#f97316", "#9333ea", "#dc2626"];

  const fetchData = () => {
    const token = localStorage.getItem("token");
    const headers = { "Authorization": `Bearer ${token}` };
    const query = filterMonth ? `?month=${filterMonth}` : "";

    fetch(`http://localhost:8080/metrics${query}`, { headers })
      .then(res => res.json())
      .then(setMetrics)
      .catch(console.error);

    fetch(`http://localhost:8080/products${query}`, { headers })
      .then(res => res.json())
      .then(setProducts)
      .catch(console.error);

    fetch(`http://localhost:8080/customers${query}`, { headers })
      .then(res => res.json())
      .then(setCustomers)
      .catch(console.error);
  };

  useEffect(() => { fetchData(); }, [filterMonth]);

  if (!metrics) return <p style={{ textAlign: "center", marginTop: 50 }}>Sem dados.</p>;

  return (
    <div style={{ padding: 20, fontFamily: "Arial, sans-serif", background: "#f4f4f4", minHeight: "100vh" }}>
      <h1 style={{ fontSize: 28, fontWeight: "bold", marginBottom: 20 }}>Dashboard Financeiro</h1>

      {/* Filtro de mês */}
      <div style={{ marginBottom: 20 }}>
        <label>Mês:</label>
        <select value={filterMonth} onChange={e => setFilterMonth(e.target.value)} style={{ marginLeft: 10 }}>
          <option value="">Todos</option>
          <option value="2025-01">Jan 2025</option>
          <option value="2025-02">Fev 2025</option>
          <option value="2025-03">Mar 2025</option>
        </select>
      </div>

      {/* Cards Totalizadores */}
      <div style={{ display: "flex", gap: 16, flexWrap: "wrap" }}>
        <div style={cardStyle}><div>Total Registros</div><div style={numberStyle}>{metrics.total_records}</div></div>
        <div style={cardStyle}><div>Total Clientes</div><div style={numberStyle}>{metrics.total_customers}</div></div>
        <div style={cardStyle}><div>Total Produtos</div><div style={numberStyle}>{metrics.total_products}</div></div>
        <div style={cardStyle}><div>Total Receita</div><div style={numberStyle}>{metrics.total_revenue.toFixed(2)}</div></div>
      </div>

      {/* Gráficos */}
      <div style={{ display: "flex", gap: 20, marginTop: 40, flexWrap: "wrap" }}>
        <div style={{ flex: 1, minWidth: 300, height: 300, background: "#fff", padding: 20, borderRadius: 8 }}>
          <h3>Top Produtos</h3>
          <ResponsiveContainer width="100%" height="80%">
            <BarChart data={products}>
              <XAxis dataKey="product" />
              <YAxis />
              <Tooltip />
              <Bar dataKey="total" fill="#2563eb" />
            </BarChart>
          </ResponsiveContainer>
        </div>

        <div style={{ flex: 1, minWidth: 300, height: 300, background: "#fff", padding: 20, borderRadius: 8 }}>
          <h3>Top Clientes</h3>
          <ResponsiveContainer width="100%" height="80%">
            <PieChart>
              <Pie
                data={customers}
                dataKey="total"
                nameKey="customer"
                cx="50%"
                cy="50%"
                outerRadius={80}
                label
              >
                {customers.map((entry, index) => (
                  <Cell key={index} fill={COLORS[index % COLORS.length]} />
                ))}
              </Pie>
              <Tooltip />
            </PieChart>
          </ResponsiveContainer>
        </div>
      </div>
    </div>
  );
}

const cardStyle = { flex: "1 1 200px", padding: 20, background: "#fff", borderRadius: 8, boxShadow: "0 2px 10px rgba(0,0,0,0.1)" };
const numberStyle = { marginTop: 10, fontSize: 22, fontWeight: "bold" };
