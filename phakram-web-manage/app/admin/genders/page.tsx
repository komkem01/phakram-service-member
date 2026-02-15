"use client";

import { useState, useEffect } from "react";
import Loading from "@/components/admin/Loading";

interface Gender {
  id: string;
  name: string;
  created_at: string;
}

export default function GendersPage() {
  const [genders, setGenders] = useState<Gender[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    // Simulate API call
    setTimeout(() => {
      setGenders([
        { id: "1", name: "ชาย", created_at: "2026-01-15" },
        { id: "2", name: "หญิง", created_at: "2026-01-15" },
        { id: "3", name: "ไม่ระบุ", created_at: "2026-01-15" },
      ]);
      setIsLoading(false);
    }, 500);
  }, []);

  if (isLoading) return <Loading />;

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-3xl font-bold text-gray-800">จัดการเพศ</h1>
        <button className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded-lg font-medium transition-colors">
          + เพิ่มเพศ
        </button>
      </div>

      <div className="bg-white rounded-lg shadow">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b border-gray-200 bg-gray-50">
                <th className="text-left py-4 px-6 text-gray-600 font-semibold">ลำดับ</th>
                <th className="text-left py-4 px-6 text-gray-600 font-semibold">ชื่อเพศ</th>
                <th className="text-left py-4 px-6 text-gray-600 font-semibold">วันที่สร้าง</th>
                <th className="text-center py-4 px-6 text-gray-600 font-semibold">จัดการ</th>
              </tr>
            </thead>
            <tbody>
              {genders.map((gender, index) => (
                <tr key={gender.id} className="border-b border-gray-100 hover:bg-gray-50">
                  <td className="py-4 px-6">{index + 1}</td>
                  <td className="py-4 px-6 font-medium">{gender.name}</td>
                  <td className="py-4 px-6 text-gray-600">{gender.created_at}</td>
                  <td className="py-4 px-6">
                    <div className="flex items-center justify-center gap-2">
                      <button className="text-blue-600 hover:text-blue-700 px-3 py-1 rounded hover:bg-blue-50 transition-colors">
                        แก้ไข
                      </button>
                      <button className="text-red-600 hover:text-red-700 px-3 py-1 rounded hover:bg-red-50 transition-colors">
                        ลบ
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
