import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import api from "../api";

export default function Home() {
  const [partners, setPartners] = useState([]);
  const [filter, setFilter] = useState("");
  const [monthFilter, setMonthFilter] = useState("");
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    fetchPartners();
  }, []);

  const fetchPartners = async () => {
    setLoading(true);
    try {
      const res = await api.get("/partners");
      setPartners(res.data || []);
    } catch (err) {
      console.error(err);
      alert("Erro ao buscar parceiros");
    } finally {
      setLoading(false);
    }
  };

  const filtered = partners.filter(p => {
    const q = filter.trim().toLowerCase();
    if (q && !(p.name || "").toLowerCase().includes(q)) return false;
    if (monthFilter) {
      if (!p.months || !p.months.includes(monthFilter)) return false;
    }
    return true;
  });

  const totalPartners = partners.length;

  return (
    <div style={styles.page}>
      <div style={styles.header}>
        <h2>Dashboard</h2>
        <div style={{ display: "flex", gap: 8 }}>
          <input placeholder="Pesquisar parceiro..." value={filter} onChange={e => setFilter(e.target.value)} style={styles.search} />
          <select value={monthFilter} onChange={e => setMonthFilter(e.target.value)} style={styles.select}>
            <option value="">Todos os meses</option>
            <option value="2025-01">Jan 2025</option>
            <option value="2025-02">Fev 2025</option>
            {/* {} */}
          </select>
          <button style={styles.primary} onClick={() => navigate("/import")}>Ir para Import</button>
        </div>
      </div>

      <div style={styles.cards}>
        <div style={styles.card}>
          <div style={styles.cardTitle}>Total de Parceiros</div>
          <div style={styles.cardNumber}>{totalPartners}</div>
        </div>

        <div style={styles.card}>
          <div style={styles.cardTitle}>Parceiros Filtrados</div>
          <div style={styles.cardNumber}>{filtered.length}</div>
        </div>
      </div>

      <div style={{ marginTop: 18 }}>
        {loading ? <div>Carregando...</div> : (
          <ul style={styles.list}>
            {filtered.map(p => (
              <li key={p.id} style={styles.item}>
                <div><strong>{p.name}</strong></div>
                <div style={{ color: "#6b7280", fontSize: 13 }}>{p.domain || p.email || ""}</div>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}

const styles = {
  page: { maxWidth: 1100, margin: "20px auto", padding: "0 16px" },
  header: { display: "flex", justifyContent: "space-between", alignItems: "center", gap: 12, marginBottom: 18 },
  search: { padding: "8px 10px", borderRadius: 8, border: "1px solid #e6e9ee" },
  select: { padding: "8px 10px", borderRadius: 8, border: "1px solid #e6e9ee" },
  primary: { padding: "8px 12px", borderRadius: 8, border: "none", background: "#2563eb", color: "#fff", cursor: "pointer" },
  cards: { display: "flex", gap: 12, marginBottom: 12 },
  card: { flex: 1, padding: 18, borderRadius: 12, background: "#fff", boxShadow: "0 6px 18px rgba(15,23,42,0.04)" },
  cardTitle: { fontSize: 13, color: "#6b7280" },
  cardNumber: { marginTop: 8, fontSize: 24, fontWeight: 800, color: "#111827" },
  list: { listStyle: "none", padding: 0, marginTop: 8 },
  item: { padding: 12, background: "#fff", borderRadius: 10, boxShadow: "0 2px 10px rgba(0,0,0,0.03)", marginBottom: 10 },
};
