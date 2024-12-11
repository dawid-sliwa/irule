import { useAuthContext } from "@/hooks/AuthContext";
import { useQuery, useQueryClient, useMutation } from "@tanstack/react-query";
import axios from "axios";
import { useForm } from "react-hook-form";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useState } from "react";

const fetchStats = async (token: string) => {
  const { data } = await axios.get("/api/v1/user-stats", {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
  return data;
};

const createUser = async (
  data: { email: string; password: string },
  token: string
) => {
  const response = await axios.post(
    "/api/v1/create-user",
    data,
    {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    }
  );
  return response.data;
};

type UserStats = {
  id: string;
  tags_created_count: number;
  email: string;
};

function Home() {
  const { token } = useAuthContext();
  const queryClient = useQueryClient();
  const { data, error, isLoading } = useQuery({
    queryKey: ["user-stats"],
    queryFn: () => fetchStats(token!),
  });

  const [isOpen, setIsOpen] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<{ email: string; password: string }>();

  const mutation = useMutation({
    mutationFn: (userData: { email: string; password: string }) =>
      createUser(userData, token!),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["user-stats"] });
      setIsOpen(false);
      reset();
    },
  });

  const onSubmit = (data: { email: string; password: string }) => {
    mutation.mutate(data);
  };

  if (isLoading)
    return <div className="text-center text-gray-400">Loading...</div>;
  if (error)
    return (
      <div className="text-center text-red-500">Error loading user stats</div>
    );

  return (
    <div className="min-h-screen flex flex-col items-center justify-start bg-gray-900 text-gray-100 py-8 px-4 relative">
      <h1 className="text-4xl font-semibold mb-8 text-center">
        User Statistics
      </h1>
      <div className="absolute top-4 right-4">
        <Dialog open={isOpen} onOpenChange={setIsOpen}>
          <DialogTrigger asChild>
            <Button variant="outline" onClick={() => setIsOpen(true)}>
              Create User
            </Button>
          </DialogTrigger>
          <DialogContent className="fixed top-20 right-4 w-80 p-6 bg-gray-800 rounded-md shadow-md">
            <DialogHeader>
              <DialogTitle className="text-xl font-bold">
                Create New User
              </DialogTitle>
            </DialogHeader>
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
              <div className="grid gap-2">
                <Label htmlFor="email">Email</Label>
                <Input
                  id="email"
                  type="email"
                  placeholder="Enter email"
                  {...register("email", { required: "Email is required" })}
                />
                {errors.email && (
                  <p className="text-red-500 text-sm">{errors.email.message}</p>
                )}
              </div>
              <div className="grid gap-2">
                <Label htmlFor="password">Password</Label>
                <Input
                  id="password"
                  type="password"
                  placeholder="Enter password"
                  {...register("password", {
                    required: "Password is required",
                  })}
                />
                {errors.password && (
                  <p className="text-red-500 text-sm">
                    {errors.password.message}
                  </p>
                )}
              </div>
              <Button type="submit" className="w-full">
                Submit
              </Button>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      <div className="max-w-4xl w-full space-y-4">
        {data?.map((user: UserStats) => (
          <div
            key={user.id}
            className="p-4 bg-gray-800 rounded-md shadow-md flex justify-between items-center"
          >
            <div>
              <p className="text-xl text-blue-500">{user.email}</p>
              <p className="text-sm text-gray-400">
                Tags Created: {user.tags_created_count}
              </p>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

export default Home;
