import { AuthContextProvider } from "@/hooks/AuthContext";
import { Outlet } from "react-router-dom";

export default function StandardLayout() {
  return (
    <div className="relative flex h-full max-w-full flex-1 flex-col overflow-hidden bg-background">
      <AuthContextProvider>
        <Outlet />
      </AuthContextProvider>
    </div>
  );
}
