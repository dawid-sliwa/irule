import { Header } from "@/components/Header";
import { useAuthContext } from "@/hooks/AuthContext";
import axios from "axios";
import { useEffect } from "react";
import { Outlet, useNavigate } from "react-router-dom";


const fetchMe = async (token: string) => {
  const response = await axios.get(
    `http://app:8080/api/v1/me`,
    {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    }
  );
  return response.data;
};

export default function Root() {
  const { isAuthenticated, setRole } = useAuthContext();
  const navigate = useNavigate();

  useEffect(() => {
    if (!isAuthenticated) {
      navigate("/login", { replace: true });
    }
  }, [isAuthenticated]);

  if (!isAuthenticated) {
    return null;
  }

  // fetch user data
  const { token } = useAuthContext();
  useEffect(() => {
    const fetchData = async (token: string) => {
      const res = await fetchMe(token!);
      setRole(res.role);
    }
    fetchData(token!);

  }, [token]);


  return (
    <div className="relative flex h-full max-w-full flex-1 flex-col overflow-hidden bg-background">
      <Header />
      <div className="p-4">
        <Outlet />
      </div>
    </div>
  );
}
