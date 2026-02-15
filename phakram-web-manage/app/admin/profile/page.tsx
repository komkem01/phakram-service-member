"use client";

import { useState, useEffect } from "react";
import Loading from "@/components/admin/Loading";

interface AdminUser {
  id: string;
  username: string;
  email: string;
  fullname: string;
  password: string;
}

export default function ProfilePage() {
  const [admin, setAdmin] = useState<AdminUser | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isEditing, setIsEditing] = useState(false);

  useEffect(() => {
    setTimeout(() => {
      setAdmin({
        id: "1",
        username: "admin",
        email: "admin@phakram.com",
        fullname: "ผู้ดูแลระบบ",
        password: "********",
      });
      setIsLoading(false);
    }, 500);
  }, []);

  if (isLoading) return <Loading />;
  if (!admin) return null;

  return (
    <div>
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-gray-800">ตั้งค่าข้อมูลผู้ดูแลระบบ</h1>
        <p className="text-gray-600 mt-2">จัดการข้อมูลส่วนตัวและการตั้งค่าบัญชี</p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Profile Card */}
        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex flex-col items-center">
            <div className="w-32 h-32 bg-gradient-to-br from-blue-500 to-purple-600 rounded-full flex items-center justify-center text-white text-4xl font-bold mb-4">
              {admin.fullname.charAt(0)}
            </div>
            <h2 className="text-xl font-bold text-gray-800">{admin.fullname}</h2>
            <p className="text-gray-600 mt-1">@{admin.username}</p>
            <p className="text-sm text-gray-500 mt-2">{admin.email}</p>
            <button className="mt-4 px-4 py-2 bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-lg transition-colors">
              เปลี่ยนรูปโปรไฟล์
            </button>
          </div>
        </div>

        {/* Information Form */}
        <div className="lg:col-span-2 bg-white rounded-lg shadow p-6">
          <div className="flex items-center justify-between mb-6">
            <h3 className="text-xl font-bold text-gray-800">ข้อมูลส่วนตัว</h3>
            <button
              onClick={() => setIsEditing(!isEditing)}
              className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-medium transition-colors"
            >
              {isEditing ? "ยกเลิก" : "แก้ไขข้อมูล"}
            </button>
          </div>

          <div className="space-y-4">
            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-2">
                ชื่อผู้ใช้งาน
              </label>
              <input
                type="text"
                value={admin.username}
                disabled={!isEditing}
                className="w-full px-4 py-3 bg-gray-50 border border-gray-200 rounded-lg text-gray-900 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:bg-gray-100 disabled:cursor-not-allowed"
              />
            </div>

            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-2">
                ชื่อ-นามสกุล
              </label>
              <input
                type="text"
                value={admin.fullname}
                disabled={!isEditing}
                className="w-full px-4 py-3 bg-gray-50 border border-gray-200 rounded-lg text-gray-900 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:bg-gray-100 disabled:cursor-not-allowed"
              />
            </div>

            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-2">อีเมล</label>
              <input
                type="email"
                value={admin.email}
                disabled={!isEditing}
                className="w-full px-4 py-3 bg-gray-50 border border-gray-200 rounded-lg text-gray-900 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:bg-gray-100 disabled:cursor-not-allowed"
              />
            </div>

            <div>
              <label className="block text-sm font-semibold text-gray-700 mb-2">
                รหัสผ่าน
              </label>
              <div className="flex gap-2">
                <input
                  type="password"
                  value={admin.password}
                  disabled={!isEditing}
                  className="flex-1 px-4 py-3 bg-gray-50 border border-gray-200 rounded-lg text-gray-900 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:bg-gray-100 disabled:cursor-not-allowed"
                />
                {isEditing && (
                  <button className="px-4 py-3 bg-gray-600 hover:bg-gray-700 text-white rounded-lg font-medium transition-colors">
                    เปลี่ยนรหัสผ่าน
                  </button>
                )}
              </div>
            </div>

            {isEditing && (
              <div className="flex gap-2 pt-4">
                <button className="flex-1 px-4 py-3 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-medium transition-colors">
                  บันทึกการเปลี่ยนแปลง
                </button>
                <button
                  onClick={() => setIsEditing(false)}
                  className="flex-1 px-4 py-3 bg-gray-200 hover:bg-gray-300 text-gray-700 rounded-lg font-medium transition-colors"
                >
                  ยกเลิก
                </button>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
