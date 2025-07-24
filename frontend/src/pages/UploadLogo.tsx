import { useState } from "react";
import axios from "axios";

// Define a type for the expected response data
interface LogoUploadResponse {
  url: string;
}

export default function UploadLogo() {
  const [file, setFile] = useState<File | null>(null);
  const [message, setMessage] = useState("");
  const [loading, setLoading] = useState(false);
  const [logoUrl, setLogoUrl] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!file) return;

    const formData = new FormData();
    formData.append("image", file);

    setLoading(true); // Start loading state

    try {
      // Send the image to the backend, use axios without explicitly typing AxiosResponse
      const response = await axios.post<LogoUploadResponse>("/api/v1/logo", formData, {
        headers: { "Content-Type": "multipart/form-data" },
      });

      // Assuming the backend returns the file URL in the response
      const uploadedLogoUrl = response.data.url;
      setLogoUrl(uploadedLogoUrl); // Save the URL to display the image
      setMessage("Logo uploaded successfully");

    } catch (error) {
      console.error("Error uploading logo:", error);
      setMessage("Failed to upload logo");
    } finally {
      setLoading(false); // Stop loading state
    }
  };

  return (
    <div className="p-6 max-w-md mx-auto bg-gray-900 text-white">
      <h2 className="text-xl font-semibold mb-4">Upload Club Logo</h2>

      <form onSubmit={handleSubmit}>
        <input
          type="file"
          onChange={(e) => setFile(e.target.files?.[0] || null)}
          className="mb-4 p-2 bg-gray-800"
        />
        <button
          className="bg-blue-600 px-4 py-2 rounded mt-4"
          type="submit"
          disabled={loading}
        >
          {loading ? "Uploading..." : "Upload"}
        </button>
      </form>

      {message && <p className="mt-2">{message}</p>}

      {/* Display the uploaded logo */}
      {logoUrl && (
        <div className="mt-4">
          <h3 className="text-lg font-semibold">Uploaded Logo:</h3>
          <img src={logoUrl} alt="Uploaded Logo" className="mt-2 max-w-xs" />
        </div>
      )}
    </div>
  );
}
