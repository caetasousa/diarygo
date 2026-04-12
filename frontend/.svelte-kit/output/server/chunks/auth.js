import { d as derived, w as writable } from "./index.js";
const BASE_URL = "/api/v1";
class ApiClient {
  token = null;
  setToken(token) {
    this.token = token;
  }
  headers(extra = {}) {
    const h = {
      "Content-Type": "application/json",
      ...extra
    };
    if (this.token) {
      h["Authorization"] = `Bearer ${this.token}`;
    }
    return h;
  }
  async request(method, path, body) {
    const res = await fetch(`${BASE_URL}${path}`, {
      method,
      headers: this.headers(),
      body: body ? JSON.stringify(body) : void 0
    });
    const data = await res.json();
    if (!res.ok) {
      const err = data;
      throw new Error(err.erro ?? `Erro ${res.status}`);
    }
    return data;
  }
  // Auth
  async registrarCliente(req) {
    return this.request("POST", "/auth/registro/cliente", req);
  }
  async registrarProfissional(req) {
    return this.request("POST", "/auth/registro/profissional", req);
  }
  async login(req) {
    return this.request("POST", "/auth/login", req);
  }
  async solicitarRecuperacao(req) {
    return this.request("POST", "/auth/solicitar-recuperacao-senha", req);
  }
  async redefinirSenha(req) {
    return this.request("POST", "/auth/redefinir-senha", req);
  }
  async me() {
    return this.request("GET", "/me");
  }
}
const api = new ApiClient();
function parsePayload(token) {
  try {
    const [, payload] = token.split(".");
    return JSON.parse(atob(payload));
  } catch {
    return null;
  }
}
function isExpired(payload) {
  return Date.now() / 1e3 > payload.exp;
}
function createAuthStore() {
  const validToken = null;
  const { subscribe, set, update } = writable({
    token: validToken,
    payload: null,
    loading: false
  });
  return {
    subscribe,
    login(token) {
      const payload = parsePayload(token);
      if (!payload || isExpired(payload)) return;
      api.setToken(token);
      set({ token, payload, loading: false });
    },
    logout() {
      api.setToken(null);
      set({ token: null, payload: null, loading: false });
    },
    setLoading(loading) {
      update((s) => ({ ...s, loading }));
    }
  };
}
const auth = createAuthStore();
const isAuthenticated = derived(auth, ($auth) => !!$auth.token);
const currentUser = derived(auth, ($auth) => $auth.payload);
export {
  currentUser as c,
  isAuthenticated as i
};
