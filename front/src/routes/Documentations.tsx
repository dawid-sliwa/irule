import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import axios from "axios";
import { useAuthContext } from "@/hooks/AuthContext";
import { Link } from "react-router-dom";
import { Button } from "@/components/ui/button";

const fetchDocumentations = async (token: string) => {
  const { data } = await axios.get("/api/v1/documentation", {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
  return data;
};

const deleteDocumentation = async ({ id, token }: { id: string; token: string }) => {
  await axios.delete(`/api/v1/documentation/${id}`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
};

type Documentation = {
  id: string;
  name: string;
  content: string;
  tag_count: number;
};

export function Documentations() {
  const { token } = useAuthContext();
  const queryClient = useQueryClient();
  const { data, error, isLoading } = useQuery({
    queryKey: ["documentation"],
    queryFn: () => fetchDocumentations(token!),
  });

  const deleteMutation = useMutation<any, Error, { id: string }>({
    mutationFn: ({ id }) => deleteDocumentation({ id, token: token! }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["documentation"] });
    },
  });
  
  if (isLoading) return <div className="text-center text-gray-400">Loading...</div>;
  if (error) return <div className="text-center text-red-500">Error loading documentations</div>;

  return (
    <div className="min-h-screen flex flex-col items-center justify-start bg-gray-900 text-gray-100 py-8 px-4">
      <h1 className="text-4xl font-semibold mb-8 text-center">Documentations</h1>
      <div className="max-w-4xl w-full space-y-4">
        {data?.map((doc: Documentation) => (
          <div key={doc.id} className="p-4 bg-gray-800 rounded-md shadow-md flex justify-between items-center">
            <div>
              <Link to={`/documentations/${doc.id}`} className="text-blue-500 hover:underline text-xl">
                {doc.name || "Untitled Documentation"}
              </Link>
              <p className="text-sm text-gray-400">Tag Count: {doc.tag_count}</p>
            </div>
            <div className="flex space-x-2">
              <Link to={`/documentations/${doc.id}/edit`} className="text-sm">
                <Button variant="outline" className="bg-green-500 text-gray-100 hover:bg-green-400">
                  Edit
                </Button>
              </Link>
              <Button
                onClick={() => deleteMutation.mutate({ id: doc.id })}
                variant="outline"
                className="bg-red-500 text-gray-100 hover:bg-red-400"
              >
                Delete
              </Button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
