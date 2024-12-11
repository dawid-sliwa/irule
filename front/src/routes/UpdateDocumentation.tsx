import { useState, useEffect } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import axios from "axios";
import { useAuthContext } from "@/hooks/AuthContext";
import { useForm, SubmitHandler } from "react-hook-form";
import { useParams, Link, useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";

interface DocumentationData {
  name: string;
  content: string;
  tags: { id: string; name: string }[];
}

interface TagData {
  id: string;
  name: string;
}

type DocumentationResponse = DocumentationData;
type UpdateDocumentationInput = Partial<DocumentationData>;

// Fetch documentation data
const fetchDocumentation = async (
  id: string,
  token: string
): Promise<DocumentationResponse> => {
  const response = await axios.get(
    `/v1/documentation/${id}`,
    {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    }
  );
  return response.data;
};

// Update documentation
const updateDocumentation = async (
  id: string,
  data: UpdateDocumentationInput,
  token: string
): Promise<DocumentationResponse> => {
  const newData = {
    title: data.name,
    content: data.content,
  }
  const response = await axios.put(
    `/v1/documentation/${id}`,
    newData,
    {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    }
  );
  return response.data;
};

// Create a new tag
const createTag = async (
  docId: string,
  tagName: string,
  token: string
): Promise<TagData> => {
  const response = await axios.post(
    `/api/v1/tag`,
    { name: tagName, documentation_id: docId },
    {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    }
  );
  return response.data;
};

// Delete a tag
const deleteTag = async (tagId: string, token: string): Promise<void> => {
  await axios.delete(`/api/v1/tag/${tagId}`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
};

export function EditDocumentation() {
  const { id } = useParams<{ id: string }>();
  const { token } = useAuthContext();
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const {
    register,
    handleSubmit,
    setValue,
    formState: { errors },
  } = useForm<DocumentationData>();

  const { data: documentationData } = useQuery({
    queryKey: ["documentation", id],
    queryFn: () => fetchDocumentation(id!, token!),
    enabled: !!id && !!token,
  });

  useEffect(() => {
    if (documentationData) {
      setValue("name", documentationData.name);
      setValue("content", documentationData.content);
    }
  }, [documentationData, setValue]);

  const updateMutation = useMutation({
    mutationFn: (data: UpdateDocumentationInput) =>
      updateDocumentation(id!, data, token!),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["documentation", id] });
      navigate(`/documentations/${id}`);
    },
  });

  const createTagMutation = useMutation({
    mutationFn: (tagName: string) => createTag(id!, tagName, token!),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["documentation", id] });
    },
  });

  const deleteTagMutation = useMutation({
    mutationFn: (tagId: string) => deleteTag(tagId, token!),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["documentation", id] });
    },
  });

  const [newTag, setNewTag] = useState<string>("");

  const onSubmit: SubmitHandler<DocumentationData> = (data) => {
    updateMutation.mutate(data);
  };

  const handleAddTag = () => {
    if (newTag) {
      createTagMutation.mutate(newTag);
      setNewTag("");
    }
  };

  const handleDeleteTag = (tagId: string) => {
    deleteTagMutation.mutate(tagId);
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-900 text-gray-100 py-8 px-4">
      <div className="max-w-lg w-full bg-gray-800 p-6 rounded-lg shadow-lg">
        <h1 className="text-3xl font-semibold mb-6 text-center">
          Edit Documentation
        </h1>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div>
            <Label className="block text-sm font-medium text-gray-300 mb-1">
              Title
            </Label>
            <Input
              id="title"
              {...register("name", {
                required: "Title is required",
                minLength: {
                  value: 5,
                  message: "Title must be at least 5 characters",
                },
              })}
              className="text-xl mt-2 text-gray-300 bg-gray-700 border-gray-600 focus:border-blue-600"
              placeholder="Enter the documentation title"
            />
            {errors.name?.message && (
              <p className="text-red-500 mt-1 text-sm">
                {String(errors.name.message)}
              </p>
            )}
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-300 mb-1">
              Content
            </label>
            <Textarea
              id="content"
              {...register("content", {
                required: "Content is required",
                minLength: {
                  value: 20,
                  message: "Content must be at least 20 characters",
                },
              })}
              rows={8}
              className="text-lg mt-2 text-gray-300 bg-gray-700 border-gray-600 focus:border-blue-600"
              placeholder="Write the documentation content here"
            />
            {errors.content && (
              <p className="text-red-500 mt-1 text-sm">
                {String(errors.content?.message)}
              </p>
            )}
          </div>
          <div className="flex space-x-2 mt-4">
            <Button
              type="submit"
              className="w-full bg-blue-500 hover:bg-blue-400 text-gray-100"
            >
              Update
            </Button>
          </div>
        </form>

        <h2 className="text-xl font-semibold mt-6 mb-4 text-center">
          Manage Tags
        </h2>
        <div className="space-y-2">
          {documentationData?.tags?.map((tag) => (
            <div key={tag.id} className="flex items-center space-x-2">
              <span className="bg-gray-700 text-gray-200 rounded-full px-3 py-1 text-sm">
                {tag.name}
              </span>
              <button
                onClick={() => handleDeleteTag(tag.id)}
                className="text-red-500 hover:text-red-400"
              >
                Delete
              </button>
            </div>
          ))}
          <div className="flex mt-4 space-x-2">
            <input
              value={newTag}
              onChange={(e) => setNewTag(e.target.value)}
              className="w-full p-2 rounded-md bg-gray-700 text-gray-100 border border-gray-600"
              placeholder="New tag name"
            />
            <Button
              onClick={handleAddTag}
              className="bg-green-500 hover:bg-green-400 text-gray-100"
            >
              Add Tag
            </Button>
          </div>
        </div>

        <div className="mt-4 text-center">
          <Link
            to={`/documentations/${id}`}
            className="text-sm underline text-gray-400 hover:text-gray-100"
          >
            Return
          </Link>
        </div>
      </div>
    </div>
  );
}
