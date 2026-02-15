"use client";

import { useState, useEffect } from "react";
import Loading from "@/components/admin/Loading";

interface Member {
  id: string;
  name: string;
  email: string;
  phone: string;
  tier: string;
  status: string;
  joined_at: string;
}

export default function MembersPage() {
  const [members, setMembers] = useState<Member[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    setTimeout(() => {
      setMembers([
        {
          id: "1",
          name: "สมชาย ใจดี",
          email: "somchai@email.com",
          phone: "081-234-5678",
          tier: "Gold",
          status: "active",
          joined_at: "2025-12-15",
        },
        {
          id: "2",
          name: "สมหญิง รักสวย",
          email: "somying@email.com",
          phone: "082-345-6789",
          tier: "Silver",
          status: "active",
          joined_at: "2026-01-10",
        },
        {
          id: "3",
          name: "วิชัย มั่นคง",
          email: "wichai@email.com",
          phone: "083-456-7890",
          tier: "Bronze",
          status: "inactive",
          joined_at: "2025-11-20",
        },
      ]);
      setIsLoading(false);
    }, 500);
  }, []);

  if (isLoading) return <Loading />;

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-3xl font-bold text-gray-800">จัดการสมาชิก</h1>
          <p className="text-gray-600 mt-2">รายการสมาชิกทั้งหมดในระบบ</p>
        </div>
        <button className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded-lg font-medium transition-colors">
          + เพิ่มสมาชิก
        </button>
      </div>

      {/* Search & Filter */}
      <div className="bg-white rounded-lg shadow p-4 mb-6">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <input
            type="text"
            placeholder="ค้นหาชื่อ, อีเมล, เบอร์โทร..."
            className="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <select className="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500">
            <option value="">ระดับสมาชิกทั้งหมด</option>
            <option value="gold">Gold</option>
            <option value="silver">Silver</option>
            <option value="bronze">Bronze</option>
          </select>
          <select className="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500">
            <option value="">สถานะทั้งหมด</option>
            <option value="active">ใช้งาน</option>
            <option value="inactive">ไม่ใช้งาน</option>
          </select>
        </div>
      </div>

      <div className="bg-white rounded-lg shadow">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b border-gray-200 bg-gray-50">
                <th className="text-left py-4 px-6 text-gray-600 font-semibold">ชื่อ-นามสกุล</th>
                <th className="text-left py-4 px-6 text-gray-600 font-semibold">อีเมล</th>
                <th className="text-left py-4 px-6 text-gray-600 font-semibold">เบอร์โทร</th>
                <th className="text-center py-4 px-6 text-gray-600 font-semibold">ระดับสมาชิก</th>
                <th className="text-center py-4 px-6 text-gray-600 font-semibold">สถานะ</th>
                <th className="text-left py-4 px-6 text-gray-600 font-semibold">วันที่สมัคร</th>
                <th className="text-center py-4 px-6 text-gray-600 font-semibold">จัดการ</th>
              </tr>
            </thead>
            <tbody>
              {members.map((member) => (
                <tr key={member.id} className="border-b border-gray-100 hover:bg-gray-50">
                  <td className="py-4 px-6 font-medium">{member.name}</td>
                  <td className="py-4 px-6 text-gray-600">{member.email}</td>
                  <td className="py-4 px-6 text-gray-600">{member.phone}</td>
                  <td className="py-4 px-6 text-center">
                    <span
                      className={`px-3 py-1 rounded-full text-sm font-medium ${
                        member.tier === "Gold"
                          ? "bg-yellow-100 text-yellow-800"
                          : member.tier === "Silver"
                          ? "bg-gray-100 text-gray-800"
                          : "bg-orange-100 text-orange-800"
                      }`}
                    >
                      {member.tier}
                    </span>
                  </td>
                  <td className="py-4 px-6 text-center">
                    <span
                      className={`px-3 py-1 rounded-full text-sm ${
                        member.status === "active"
                          ? "bg-green-100 text-green-800"
                          : "bg-red-100 text-red-800"
                      }`}
                    >
                      {member.status === "active" ? "ใช้งาน" : "ไม่ใช้งาน"}
                    </span>
                  </td>
                  <td className="py-4 px-6 text-gray-600">{member.joined_at}</td>
                  <td className="py-4 px-6">
                    <div className="flex items-center justify-center gap-2">
                      <button className="text-blue-600 hover:text-blue-700 px-3 py-1 rounded hover:bg-blue-50 transition-colors">
                        ดู
                      </button>
                      <button className="text-green-600 hover:text-green-700 px-3 py-1 rounded hover:bg-green-50 transition-colors">
                        แก้ไข
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
