import { Link, useNavigate } from "react-router-dom";
import { useAuthContext } from "@/hooks/AuthContext";
import { Button } from "@/components/ui/button";

export function Header() {
  const { logout, role } = useAuthContext();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate("/login");
  };

  return (
    <header className="w-full p-4 bg-gray-800 flex justify-between items-center text-gray-100 shadow-lg">
      <nav className="flex space-x-4">
        {role === "admin" && (
          <Link to="/" className="text-lg hover:text-blue-400">Home</Link>
        )}
        <Link to="/documentations" className="text-lg hover:text-blue-400">Documentations</Link>
        <Link to="/documentations/create" className="text-lg hover:text-blue-400">Create Documentation</Link>
      </nav>
      <Button onClick={handleLogout} className="bg-red-600 hover:bg-red-700">
        Logout
      </Button>
    </header>
  );
}
