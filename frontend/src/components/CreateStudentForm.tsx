import type React from "react";
import { useState } from "react";
import api from "../api";
import type { AxiosError } from "axios";

const CreateStudentForm: React.FC = () => {
    const [name, setName] = useState('');
    const [beltGrade, setBeltGrade] = useState('');
    const [dob, setDob] = useState('');
    const [studentType, setStudentType] = useState<'kid' | 'adult'>('kid');
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
                studentType,
            });

            setSuccess('Student created successfully!');
            setName('');
            setBeltGrade('');
            setDob('');
            setStudentType('kid');
        } catch (err) {
            const error = err as AxiosError<{ error: string }>;
            setError(error.response?.data?.error || "Failed to create student");
            console.error("Error creating student:", error);
        }
    };

    const formatDateForBackend = (dateStr: string): string => {
        const [year, month, day] = dateStr.split("-");
        return `${day}/${month}/${year}`;
    };

    return (
        <div className="max-w-md mx-auto mt-10 bg-gray-900 text-white p-6 rounded-lg shadow">
            <h2 className="text-2xl font-bold mb-6 text-center">Create Student</h2>
            <form onSubmit={handleSubmit} className="space-y-5">
                <div>
                    <label className="block font-medium mb-1">Name</label>
                    <input
                        type="text"
                        className="w-full bg-gray-800 border border-gray-700 px-3 py-2 rounded text-white focus:outline-none focus:border-blue-500"
                        value={name}
                        onChange={(e) => setName(e.target.value)}
                        required
                    />
                </div>

                <div>
                    <label className="block font-medium mb-1">Belt Grade</label>
                    <input
                        type="text"
                        className="w-full bg-gray-800 border border-gray-700 px-3 py-2 rounded text-white focus:outline-none focus:border-blue-500"
                        value={beltGrade}
                        onChange={(e) => setBeltGrade(e.target.value)}
                        required
                    />
                </div>

                <div>
                    <label className="block font-medium mb-1">Date of Birth</label>
                    <input
                        type="date"
                        className="w-full bg-gray-800 border border-gray-700 px-3 py-2 rounded text-white focus:outline-none focus:border-blue-500"
                        value={dob}
                        onChange={(e) => setDob(e.target.value)}
                        required
                    />
                </div>

                <div>
                    <label className="block font-medium mb-1">Student Type</label>
                    <div className="flex gap-4">
                        <label className="flex items-center">
                            <input
                                type="radio"
                                name="studentType"
                                value="kid"
                                checked={studentType === 'kid'}
                                onChange={() => setStudentType('kid')}
                                className="mr-2 accent-blue-600"
                            />
                            Kid
                        </label>
                        <label className="flex items-center">
                            <input
                                type="radio"
                                name="studentType"
                                value="adult"
                                checked={studentType === 'adult'}
                                onChange={() => setStudentType('adult')}
                                className="mr-2 accent-blue-600"
                            />
                            Adult
                        </label>
                    </div>
                </div>

                <button
                    type="submit"
                    className="w-full bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-400"
                >
                    Create Student
                </button>

                {success && <p className="text-green-400 text-sm mt-2">{success}</p>}
                {error && <p className="text-red-400 text-sm mt-2">{error}</p>}
            </form>
        </div>
    );
};

export default CreateStudentForm;
