import { RouterProvider } from "react-router-dom";
import { Router } from "./routes";
import { ThemeProvider } from "./hooks/ThemeContext";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

function App() {
  const queryClient = new QueryClient();

  return (
    <QueryClientProvider client={queryClient}>
      <ThemeProvider defaultTheme="dark">
        <RouterProvider router={Router} />
      </ThemeProvider>
    </QueryClientProvider>
  );
}

export default App;
