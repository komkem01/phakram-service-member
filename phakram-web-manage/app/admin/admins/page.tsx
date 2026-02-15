"use client";

import { useState, useEffect } from "react";
import Loading from "@/components/admin/Loading";

export default function AdminsPage() {
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    setTimeout(() => setIsLoading(false), 500);
  }, []);

  if (isLoading) return <Loading />;

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-3xl font-bold text-gray-800">จัดการผู้ดูแลระบบ</h1>
          <p className="text-gray-600 mt-2">รายการผู้ดูแลระบบทั้งหมด</p>
        </div>
        <button className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded-lg font-medium transition-colors">
          + เพิ่มผู้ดูแลระบบ
        </button>
      </div>

      <div className="bg-white rounded-lg shadow">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b border-gray-200 bg-gray-50">
                <th className="text-left py-4 px-6 text-gray-600 font-semibold">ชื่อผู้ใช้</th>
                <th className="text-left py-4 px-6 text-gray-600 font-semibold">ชื่อ-นามสกุล</th>
                <th className="text-left py-4 px-6 text-gray-600 font-semibold">อีเมล</th>
                <th className="text-center py-4 px-6 text-gray-600 font-semibold">สถานะ</th>
                <th className="text-left py-4 px-6 text-gray-600 font-semibold">วันที่สร้าง</th>
                <th className="text-center py-4 px-6 text-gray-600 font-semibold">จัดการ</th>
              </tr>
            </thead>
            <tbody>
              <tr className="border-b border-gray-100 hover:bg-gray-50">
                <td className="py-4 px-6 font-medium">admin</td>
                <td className="py-4 px-6">ผู้ดูแลระบบ</td>
                <td className="py-4 px-6 text-gray-600">admin@phakram.com</td>
                <td className="py-4 px-6 text-center">
                  <span className="bg-green-100 text-green-800 px-3 py-1 rounded-full text-sm">
                    ใช้งาน
                  </span>
                </td>
                <td className="py-4 px-6 text-gray-600">2026-01-01</td>
                <td className="py-4 px-6">
                  <div className="flex items-center justify-center gap-2">
                    <button className="text-blue-600 hover:text-blue-700 px-3 py-1 rounded hover:bg-blue-50 transition-colors">
                      แก้ไข
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
