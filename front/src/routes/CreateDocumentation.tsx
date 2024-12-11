import { useAuthContext } from "@/hooks/AuthContext";
import { useQueryClient, useMutation } from "@tanstack/react-query";
import axios from "axios";
import { useForm } from "react-hook-form";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Dialog, DialogContent, DialogTitle, DialogDescription, DialogClose } from "@/components/ui/dialog";
import { useState } from "react";

const createDocumentation = async (data: any, token: string) => {
  const response = await axios.post(
    "http://app:8080/api/v1/documentation",
    data,
    {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    }
  );
  return response.data;
};

export function CreateDocumentation() {
  const { token } = useAuthContext();
  const queryClient = useQueryClient();
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm();
  const [open, setOpen] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const mutation = useMutation({
    mutationFn: (data: any) => createDocumentation(data, token!),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["documentation"] });
      setOpen(true);
      reset();
    },
    onError: (error: any) => {
      if (error.response?.status === 400) {
        setError("Invalid input. Please check the fields and try again.");
      } else {
        setError("An unexpected error occurred. Please try again.");
      }
    },
  });

  const onSubmit = (data: any) => {
    mutation.mutate(data);
  };

  return (
    <div className="flex h-screen w-full items-center justify-center bg-gray-900 px-4">
      <Card className="w-full max-w-lg bg-gray-800 shadow-xl rounded-lg">
        <CardHeader className="text-center">
          <CardTitle className="text-2xl font-semibold text-gray-100">
            Create Documentation
          </CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
            <div>
              <Label htmlFor="title" className="text-lg text-gray-300">
                Title
              </Label>
              <Input
                id="title"
                {...register("title", {
                  required: "Title is required",
                  minLength: { value: 5, message: "Title must be at least 5 characters" },
                })}
                className="text-xl mt-2 text-gray-300 bg-gray-700 border-gray-600 focus:border-blue-600"
                placeholder="Enter the documentation title"
              />
              {errors.title?.message && (
                <p className="text-red-500 mt-1 text-sm">{String(errors.title.message)}</p>
              )}
            </div>
            <div>
              <Label htmlFor="content" className="text-lg text-gray-300">
                Content
              </Label>
              <Textarea
                id="content"
                {...register("content", {
                  required: "Content is required",
                  minLength: { value: 20, message: "Content must be at least 20 characters" },
                })}
                rows={8}
                className="text-lg mt-2 text-gray-300 bg-gray-700 border-gray-600 focus:border-blue-600"
                placeholder="Write the documentation content here"
              />
              {errors.content && (
                <p className="text-red-500 mt-1 text-sm">{String(errors.content?.message)}</p>
              )}
            </div>
            <div className="flex items-center justify-center">
              <Button
                type="submit"
                className="bg-blue-600 hover:bg-blue-700 text-white font-semibold text-lg px-8 py-3 rounded-md focus:ring-4 focus:ring-blue-500 focus:ring-opacity-50"
              >
                Create Documentation
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>

      {/* Dialog for Success Confirmation */}
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent className="bg-gray-800 text-gray-100 rounded-lg shadow-lg">
          <DialogTitle className="text-2xl font-semibold">Success</DialogTitle>
          <DialogDescription className="text-lg mt-2">
            The documentation has been created successfully.
          </DialogDescription>
          <DialogClose asChild>
            <Button className="mt-6 bg-blue-600 hover:bg-blue-700 text-white font-semibold py-2 px-4 rounded-md">
              Close
            </Button>
          </DialogClose>
        </DialogContent>
      </Dialog>

      {/* Dialog for Error Message */}
      <Dialog open={!!error} onOpenChange={() => setError(null)}>
        <DialogContent className="bg-gray-800 text-gray-100 rounded-lg shadow-lg">
          <DialogTitle className="text-2xl font-semibold">Error</DialogTitle>
          <DialogDescription className="text-lg mt-2">
            {error}
          </DialogDescription>
          <DialogClose asChild>
            <Button className="mt-6 bg-red-600 hover:bg-red-700 text-white font-semibold py-2 px-4 rounded-md">
              Close
            </Button>
          </DialogClose>
        </DialogContent>
      </Dialog>
    </div>
  );
}
