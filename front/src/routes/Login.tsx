import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { useForm } from "react-hook-form";
import { TloginUser } from "@/common/types";
import { useAuthContext } from "@/hooks/AuthContext";
import { useMutation } from "@tanstack/react-query";
import axios from "axios";

export function Login() {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<TloginUser>();
  
  const { login } = useAuthContext();

  const mutation = useMutation({
    mutationFn: (data: TloginUser) => axios.post("/api/v1/auth/login", data),
    onSuccess: (response) => {
      const token = response.data.token;
      login(token);
    },
    onError: (error) => {
      console.error("An error occurred:", error);
    },
  });
  const onSubmit = (data: TloginUser) => {
    mutation.mutate(data);
  };

  return (
    <div className="flex h-screen w-full items-center justify-center px-4 text-foreground">
      <Card className="mx-auto max-w-sm">
        <CardHeader>
          <CardTitle className="text-2xl">Login</CardTitle>
          <CardDescription>
            Enter your email below to login to your account
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit(onSubmit)}>
            <div className="grid gap-4">
              <div className="grid gap-2">
                <Label htmlFor="email">Email</Label>
                <Input
                  id="email"
                  type="email"
                  placeholder="m@example.com"
                  {...register("email")}
                />
                {errors.email && (
                  <p className="text-red-500 text-sm">{errors.email.message}</p>
                )}
              </div>
              <div className="grid gap-2">
                <div className="flex items-center">
                  <Label htmlFor="password">Password</Label>
                  <a
                    href="#"
                    className="ml-auto inline-block text-sm underline"
                  >
                    Forgot your password?
                  </a>
                </div>
                <Input
                  id="password"
                  type="password"
                  {...register("password")}
                />
                {errors.password && (
                  <p className="text-red-500 text-sm">
                    {errors.password.message}
                  </p>
                )}
              </div>
              <Button type="submit" className="w-full">
                Login
              </Button>
            </div>
            <div className="mt-4 text-center text-sm">
              Don&apos;t have an account?{" "}
              <a href="/register" className="underline">
                Sign up
              </a>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
