import { useState, useEffect } from "react";
import axios from "axios";
import type {
    GroupedAccounts,
    GroupedMonthAccounts,
    GroupedDayAccounts,
    AccountEntry,
    AccountRow,
} from "../types/Accounts";
/* eslint-disable  @typescript-eslint/no-explicit-any */
type EntryField = "description" | "amount";

export default function AccountsPage() {
    const [income, setIncome] = useState<AccountRow[]>([{ description: "", amount: "" }]);
    const [expenditure, setExpenditure] = useState<AccountRow[]>([{ description: "", amount: "" }]);
    const [accounts, setAccounts] = useState<GroupedAccounts[]>([]);
    const [error, setError] = useState<string>("");
    const [success, setSuccess] = useState<string>("");
    const [view, setView] = useState<"form" | "list">("form");

    const handleAddRow = (type: "income" | "expenditure") => {
        const newRow: AccountRow = { description: "", amount: "" };
        if (type === "income") setIncome([...income, newRow]);
        else setExpenditure([...expenditure, newRow]);
    };

    const handleChange = (
        type: "income" | "expenditure",
        index: number,
        field: EntryField,
        value: string
    ) => {
        const rows = type === "income" ? [...income] : [...expenditure];
        rows[index][field] = value;
        if (type === "income") setIncome(rows);
        else setExpenditure(rows);
    };

    const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setSuccess("");

    // Remove invalid/empty rows and convert to numbers
    const cleanedIncome = income
        .filter(i => i.description.trim() && !isNaN(parseFloat(i.amount)) && parseFloat(i.amount) > 0)
        .map(i => ({ ...i, amount: parseFloat(i.amount) }));

    const cleanedExpenditure = expenditure
        .filter(e => e.description.trim() && !isNaN(parseFloat(e.amount)) && parseFloat(e.amount) > 0)
        .map(e => ({ ...e, amount: parseFloat(e.amount) }));

    if (cleanedIncome.length === 0 && cleanedExpenditure.length === 0) {
        setError("Please add at least one valid income or expenditure entry.");
        return;
    }

    try {
        await axios.post("/api/v1/accounts", {
            income: cleanedIncome,
            expenditure: cleanedExpenditure,
        });
        setIncome([{ description: "", amount: "" }]);
        setExpenditure([{ description: "", amount: "" }]);
        setSuccess("Accounts added successfully!");
        fetchAccounts();
        setView("list");
        setTimeout(() => setSuccess(""), 3000);
    } catch (err: any) {
        console.error("Save failed:", err.response?.data || err.message);
        setError(err.response?.data?.error || "Failed to save accounts.");
    }
};


    const fetchAccounts = async () => {
        const res = await axios.get<GroupedAccounts[]>("/api/v1/accounts");
        setAccounts(res.data);
    };

    const handleRemoveRow = (type: "income" | "expenditure", index: number) => {
        if (type === "income") {
            setIncome(income.filter((_, i) => i !== index));
        } else {
            setExpenditure(expenditure.filter((_, i) => i !== index));
        }
    };

    useEffect(() => {
        if (view === "list") fetchAccounts();
    }, [view]);

    return (
        <div className="p-6 max-w-4xl mx-auto text-white bg-gray-900 min-h-screen">
            <div className="flex justify-between items-center mb-6">
                <h1 className="text-2xl font-bold">Club Accounts</h1>
                <button
                    onClick={() => setView(view === "form" ? "list" : "form")}
                    className="bg-gray-700 px-4 py-2 rounded"
                >
                    {view === "form" ? "View Accounts Overview" : "Add New Accounts"}
                </button>
            </div>

            {success && <div className="bg-green-600 p-3 rounded mb-4">{success}</div>}
            {error && <div className="bg-red-600 p-3 rounded mb-4">{error}</div>}

            {view === "form" && (
                <form onSubmit={handleSubmit} className="space-y-6 mb-10">
                    <div>
                        <h2 className="text-xl mb-2">Incomes</h2>
                        {income.map((row, idx) => (
                            <div key={idx} className="flex gap-2 mb-2 items-center">
                                <input
                                    className="p-2 rounded bg-gray-800 w-1/2"
                                    placeholder="Description"
                                    value={row.description}
                                    onChange={(e) => handleChange("income", idx, "description", e.target.value)}
                                />
                                <input
                                    className="p-2 rounded bg-gray-800 w-1/4"
                                    placeholder="Amount"
                                    type="number"
                                    step="0.01"
                                    value={row.amount}
                                    onChange={(e) => handleChange("income", idx, "amount", e.target.value)}
                                />
                                {income.length > 1 && (
                                    <button
                                        type="button"
                                        onClick={() => handleRemoveRow("income", idx)}
                                        className="bg-red-600 px-2 py-1 rounded"
                                    >
                                        X
                                    </button>
                                )}
                            </div>
                        ))}
                        <button
                            type="button"
                            className="bg-blue-600 px-3 py-1 rounded mt-1"
                            onClick={() => handleAddRow("income")}
                        >
                            + Add Income
                        </button>
                    </div>

                    <div>
                        <h2 className="text-xl mb-2">Expenditures</h2>
                        {expenditure.map((row, idx) => (
                            <div key={idx} className="flex gap-2 mb-2">
                                <input
                                    className="p-2 rounded bg-gray-800 w-1/2"
                                    placeholder="Description"
                                    value={row.description}
                                    onChange={(e) => handleChange("expenditure", idx, "description", e.target.value)}
                                />
                                <input
                                    className="p-2 rounded bg-gray-800 w-1/4"
                                    placeholder="Amount"
                                    type="number"
                                    step="0.01"
                                    value={row.amount}
                                    onChange={(e) => handleChange("expenditure", idx, "amount", e.target.value)}
                                />

                                {expenditure.length > 1 && (
                                    <button
                                        type="button"
                                        onClick={() => handleRemoveRow("income", idx)}
                                        className="bg-red-600 px-2 py-1 rounded"
                                    >
                                        X
                                    </button>
                                )}
                            </div>
                        ))}
                        <button
                            type="button"
                            className="bg-blue-600 px-3 py-1 rounded mt-1"
                            onClick={() => handleAddRow("expenditure")}
                        >
                            + Add Expenditure
                        </button>
                    </div>

                    <button type="submit" className="bg-green-600 px-4 py-2 rounded">
                        Save Accounts
                    </button>
                </form>
            )}

            {view === "list" && (
                <div>
                    <h2 className="text-xl mb-4">Accounts Overview</h2>
                    {accounts.length > 0 ? (
                        accounts.map((year: GroupedAccounts) => (
                            <div key={year.year} className="mb-6">
                                <h3 className="text-lg font-semibold">{year.year}</h3>
                                {year.months.map((month: GroupedMonthAccounts) => (
                                    <div key={month.month} className="ml-4 mb-4">
                                        <h4 className="text-md font-semibold">{month.month}</h4>
                                        {month.entries.map((day: GroupedDayAccounts) => (
                                            <div key={day.date} className="ml-6 mb-2">
                                                <p className="font-medium">{day.date}</p>
                                                <ul className="ml-4">
                                                    {day.records.map((record: AccountEntry, idx: number) => (
                                                        <li key={idx}>
                                                            {record.description} - £{record.amount.toFixed(2)} ({record.category})
                                                        </li>
                                                    ))}
                                                </ul>
                                            </div>
                                        ))}
                                    </div>
                                ))}
                            </div>
                        ))
                    ) : (
                        <p>No accounts yet.</p>
                    )}
                </div>
            )}
        </div>
    );
}
