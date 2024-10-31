import { createBrowserRouter, Outlet } from "react-router-dom";
import Home from "./Home";
import { Login } from "./Login";
import Root from "./Root";
import { Register } from "./Register";
import { AuthContextProvider } from "@/hooks/AuthContext";
import StandardLayout from "./Layouts/StandardLayout";
import { Documentations } from "./Documentations";
import { Documentation } from "./Documentation";
import { CreateDocumentation } from "./CreateDocumentation";
import { EditDocumentation } from "./UpdateDocumentation";

const AuthLayout = () => (
  <AuthContextProvider>
    <Outlet />
  </AuthContextProvider>
);

export const Router = createBrowserRouter([
  {
    path: "/",
    element: <AuthLayout />,
    children: [
      {
        path: "/",
        element: <Root />,
        children: [
          {
            path: "/",
            element: <Home />,
          },
          {
            path: "/documentations",
            element: <Documentations />,
          },
          {
            path: "/documentations/:id",
            element: <Documentation />,
          },
          {
            path: "/documentations/create",
            element: <CreateDocumentation />,
          },
          {
            path: "/documentations/:id/edit",
            element: <EditDocumentation />,
          },
        ],
      },
    ],
  },
  {
    path: "/",
    element: <StandardLayout />,
    children: [
      {
        path: "/login",
        element: <Login />,
      },
      {
        path: "/register",
        element: <Register />,
      },
    ],
  },
]);
