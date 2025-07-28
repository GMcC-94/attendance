export interface AccountEntry {
  description: string;
  amount: number;
  category: "income" | "expenditure";
}

export interface GroupedDayAccounts {
  date: string;
  records: AccountEntry[];
}

export interface GroupedMonthAccounts {
  month: string;
  entries: GroupedDayAccounts[];
}

export interface GroupedAccounts {
  year: number;
  months: GroupedMonthAccounts[];
}

export interface AccountRow {
  description: string;
  amount: string;
}