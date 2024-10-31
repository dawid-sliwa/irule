import { useQuery } from "@tanstack/react-query";
import axios from "axios";
import { useAuthContext } from "@/hooks/AuthContext";
import { useParams, Link } from "react-router-dom";
import { Button } from "@/components/ui/button";

const fetchDocumentation = async (id: string, token: string) => {
  const { data } = await axios.get(`http://localhost:8080/api/v1/documentation/${id}`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
  return data;
};

export function Documentation() {
  const { id } = useParams<{ id: string }>();
  const { token } = useAuthContext();
  const { data, error, isLoading } = useQuery({
    queryKey: ["documentation", id],
    queryFn: () => fetchDocumentation(id!, token!),
  });

  if (isLoading) return <div className="text-center text-gray-400">Loading...</div>;
  if (error) return <div className="text-center text-red-500">Error loading documentation</div>;

  return (
    <div className="min-h-screen flex flex-col items-center justify-start bg-gray-900 text-gray-100 py-8 px-4">
      <div className="max-w-3xl w-full p-6 bg-gray-800 rounded-lg shadow-lg">
        <div className="flex justify-between items-center mb-4">
          <h1 className="text-3xl font-semibold">
            {data.name || "Untitled Documentation"}
          </h1>
          <Link to={`/documentations/${id}/edit`}>
            <Button variant="outline" className="bg-blue-500 text-gray-100 hover:bg-blue-400">
              Edit
            </Button>
          </Link>
        </div>
        <p className="text-lg leading-relaxed mb-6">{data.content || "No content available."}</p>
        <p className="text-sm text-gray-400 mb-2">Tag Count: {data.tag_count}</p>
        {data.tags && (
          <ul className="flex flex-wrap gap-2">
            {data.tags.map((tag: { id: string; name: string }) => (
              <li
                key={tag.id}
                className="bg-gray-700 px-3 py-1 rounded-full text-sm text-gray-300"
              >
                {tag.name}
              </li>
            ))}
          </ul>
        )}
        <div className="mt-6">
          <Link to="/documentations" className="text-sm underline text-gray-400 hover:text-gray-100">
            Back to Documentations
          </Link>
        </div>
      </div>
    </div>
  );
}
