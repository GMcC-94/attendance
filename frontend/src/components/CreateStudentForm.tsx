import type React from "react";
import { useState } from "react";
import api from "../api";

/* eslint-disable  @typescript-eslint/no-explicit-any */
const CreateStudentForm: React.FC = () => {
    const [name, setName] = useState('');
    const [beltGrade, setBeltGrade] = useState('');
    const [dob, setDob] = useState('');
    const [success, setSuccess] = useState<string | null>(null);
    const [error, setError] = useState<string | null>(null);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setSuccess(null);
        setError(null);

        try {
            await api.post('/students', {
                name,
                beltGrade,
                dob: formatDateForBackend(dob),
            });

            setSuccess('Student created successfully!')
            setName('');
            setBeltGrade('');
            setDob('');
        } catch (err: any) {
            setError(err.response?.data || 'Failed to create student')
            console.error('Error creating student:', err)
        }
    };

    const formatDateForBackend = (dateStr: string): string => {
        const [year, month, day] = dateStr.split("-");
        return `${day}/${month}/${year}`;
    };


    return (
        <div className="max-w-md mx-auto mt-10 bg-white p-6 rounded shadow">
            <h2 className="text-2xl font-bold mb-4 text-center">Create Student</h2>
            <form onSubmit={handleSubmit} className="space-y-4">
                <div>
                    <label className="block font-medium">Name</label>
                    <input
                        type="text"
                        className="w-full border px-3 py-2 rounded"
                        value={name}
                        onChange={(e) => setName(e.target.value)}
                        required
                    />
                </div>

                <div>
                    <label className="block font-medium">Belt Grade</label>
                    <input
                        type="text"
                        className="w-full border px-3 py-2 rounded"
                        value={beltGrade}
                        onChange={(e) => setBeltGrade(e.target.value)}
                        required
                    />
                </div>

                <div>
                    <label className="block font-medium">Date of Birth</label>
                    <input
                        type="date"
                        className="w-full border px-3 py-2 rounded"
                        value={dob}
                        onChange={(e) => setDob(e.target.value)}
                        required
                    />
                </div>

                <button
                    type="submit"
                    className="w-full bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
                >
                    Create Student
                </button>

                {success && <p className="text-green-600 text-sm mt-2">{success}</p>}
                {error && <p className="text-red-600 text-sm mt-2">{error}</p>}
            </form>
        </div>
    );
};

export default CreateStudentForm;