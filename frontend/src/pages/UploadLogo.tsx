import { useState } from "react";
import axios from "axios";

export default function UploadLogo() {
  const [file, setFile] = useState<File | null>(null);
  const [message, setMessage] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!file) return;

    const formData = new FormData();
    formData.append("image", file);

    try {
      await axios.post("/api/v1/logo", formData, {
        headers: { "Content-Type": "multipart/form-data" },
      });
      setMessage("Logo uploaded successfully");
    } catch {
      setMessage("Failed to upload logo");
    }
  };

  return (
    <div className="p-6 max-w-md mx-auto bg-gray-900 text-white">
      <h2 className="text-xl font-semibold mb-4">Upload Club Logo</h2>
      <form onSubmit={handleSubmit}>
        <input type="file" onChange={(e) => setFile(e.target.files?.[0] || null)} />
        <button className="bg-blue-600 px-4 py-2 rounded mt-4" type="submit">
          Upload
        </button>
      </form>
      <p className="mt-2">{message}</p>
    </div>
  );
}