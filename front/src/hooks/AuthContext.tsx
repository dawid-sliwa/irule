import {
  useState,
  useCallback,
  createContext,
  useContext,
  ReactNode,
} from "react";
import { useNavigate } from "react-router-dom";

const AuthContext = createContext<AuthContextType | undefined>(undefined);

type AuthContextType = {
  token: string | null;
  isAuthenticated: boolean;
  login: (token: string) => void;
  logout: () => void;
  role: string | null;
  setRole: (role: string) => void;
};

const AuthContextProvider = ({ children }: { children: ReactNode }) => {
  const [token, setToken] = useState<string | null>(
    localStorage.getItem("token")
  );

  const [role, setRole] = useState<string | null>(
    "user"
  );

  const navigate = useNavigate();

  const isAuthenticated = !!token;

  const login = useCallback(
    (token: string) => {
      localStorage.setItem("token", token);
      setToken(token);
      navigate("/", { replace: true });
    },
    [navigate]
  );

  const logout = useCallback(() => {
    localStorage.removeItem("token");
    setToken(null);
    navigate("/login", { replace: true });
  }, [navigate]);

  return (
    <AuthContext.Provider value={{ token, isAuthenticated, login, logout, role, setRole }}>
      {children}
    </AuthContext.Provider>
  );
};

const useAuthContext = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error(
      "useAuthContext must be used within an AuthContextProvider"
    );
  }
  return context;
};

export { AuthContextProvider, useAuthContext };
