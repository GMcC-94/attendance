import { useState } from "react";
import axios from "axios";

export default function UploadLogo() {
  const [message, setMessage] = useState("");

  const handleChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    const formData = new FormData();
    formData.append("logo", file);

    try {
      await axios.post("/api/v1/logo", formData, {
        headers: { "Content-Type": "multipart/form-data" },
      });
      setMessage("Logo uploaded!");
    } catch {
      setMessage("Failed to upload logo");
    }
  };

  return (
    <div className="p-4">
      <h2 className="text-xl font-bold mb-2">Upload Club Logo</h2>
      <input type="file" accept="image/*" onChange={handleChange} />
      {message && <p className="mt-2">{message}</p>}
    </div>
  );
}
