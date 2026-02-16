"use client";

import { useState, useEffect, useCallback, useRef } from "react";
import { District, districtsApi } from "@/lib/api/districts";
import { provincesApi } from "@/lib/api/provinces";
import { formatDateOnly } from "@/lib/utils/date";
import Loading from "@/components/admin/Loading";
import DistrictFormModal from "@/components/admin/DistrictFormModal";
import { useModal } from "@/contexts/ModalContext";

export default function DistrictsPage() {
  const [districts, setDistricts] = useState<District[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [total, setTotal] = useState(0);
  const [provinceNameMap, setProvinceNameMap] = useState<Record<string, string>>({});
  const provinceNameMapRef = useRef<Record<string, string>>({});
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [modalMode, setModalMode] = useState<'create' | 'edit'>('create');
  const [selectedDistrict, setSelectedDistrict] = useState<District | null>(null);

  const { showSuccess, showError, showWarning } = useModal();

  const totalPages = Math.max(1, Math.ceil(total / pageSize));

  useEffect(() => {
    provinceNameMapRef.current = provinceNameMap;
  }, [provinceNameMap]);

  const hydrateMissingProvinceNames = useCallback(async (items: District[]) => {
    const neededIds = Array.from(new Set(items.map((item) => item.province_id).filter(Boolean)));
    const missingIds = neededIds.filter((id) => !provinceNameMapRef.current[id]);
    if (missingIds.length === 0) return;

    const results = await Promise.allSettled(
      missingIds.map(async (id) => {
        const province = await provincesApi.getById(id);
        return { id, name: province.name };
      })
    );

    const fallbackMap = results.reduce<Record<string, string>>((acc, result) => {
      if (result.status === 'fulfilled') {
        acc[result.value.id] = result.value.name;
      }
      return acc;
    }, {});

    if (Object.keys(fallbackMap).length > 0) {
      setProvinceNameMap((prev) => ({ ...prev, ...fallbackMap }));
    }
  }, []);

  const fetchProvinces = useCallback(async () => {
    try {
      const response = await provincesApi.list({ page: 1, size: 1000 });
      const mapping = response.data.reduce<Record<string, string>>((acc, province) => {
        acc[province.id] = province.name;
        return acc;
      }, {});
      setProvinceNameMap(mapping);
    } catch (error) {
      console.error('Failed to load provinces:', error);
    }
  }, []);

  const fetchDistricts = useCallback(async (currentPage: number, currentSize: number) => {
    try {
      setIsLoading(true);
      const response = await districtsApi.list({ page: currentPage, size: currentSize });
      setDistricts(response.data);
      setTotal(response.paginate.total);
      setPage(response.paginate.page);
      await hydrateMissingProvinceNames(response.data);
    } catch {
      showError('ไม่สามารถโหลดข้อมูลอำเภอได้');
    } finally {
      setIsLoading(false);
    }
  }, [hydrateMissingProvinceNames, showError]);

  useEffect(() => {
    fetchDistricts(page, pageSize);
  }, [fetchDistricts, page, pageSize]);

  useEffect(() => {
    fetchProvinces();
  }, [fetchProvinces]);

  const handleCreate = () => {
    setModalMode('create');
    setSelectedDistrict(null);
    setIsModalOpen(true);
  };

  const handleEdit = (district: District) => {
    setModalMode('edit');
    setSelectedDistrict(district);
    setIsModalOpen(true);
  };

  const handleDelete = (district: District) => {
    showWarning(
      `ต้องการลบอำเภอ "${district.name}" หรือไม่?`,
      async () => {
        try {
          await districtsApi.delete(district.id);
          showSuccess('ลบอำเภอสำเร็จ');
          fetchDistricts(page, pageSize);
        } catch {
          showError('ไม่สามารถลบอำเภอได้');
        }
      },
      'ยืนยันการลบ'
    );
  };

  const handleSubmit = async (data: { province_id: string; name: string; is_active: boolean }) => {
    try {
      if (modalMode === 'create') {
        await districtsApi.create(data);
        showSuccess('เพิ่มอำเภอสำเร็จ');
      } else if (selectedDistrict) {
        await districtsApi.update(selectedDistrict.id, data);
        showSuccess('แก้ไขอำเภอสำเร็จ');
      }
      setIsModalOpen(false);
      fetchDistricts(page, pageSize);
    } catch {
      showError(`ไม่สามารถ${modalMode === 'create' ? 'เพิ่ม' : 'แก้ไข'}อำเภอได้`);
    }
  };

  const handlePageChange = (newPage: number) => {
    if (newPage < 1 || newPage > totalPages) return;
    setPage(newPage);
  };

  const handlePageSizeChange = (newSize: number) => {
    setPageSize(newSize);
    setPage(1);
  };

  if (isLoading) return <Loading />;

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-3xl font-bold text-gray-800">จัดการอำเภอ</h1>
        <button
          onClick={handleCreate}
          className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded-lg font-medium transition-colors"
        >
          + เพิ่มอำเภอ
        </button>
      </div>

      <div className="bg-white rounded-lg shadow overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                ชื่ออำเภอ
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                จังหวัด
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                สถานะ
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                วันที่สร้าง
              </th>
              <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                จัดการ
              </th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {districts.map((district) => (
              <tr key={district.id}>
                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                  {district.name}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-700">
                  {provinceNameMap[district.province_id] ?? '-'}
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span
                    className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                      district.is_active
                        ? 'bg-green-100 text-green-800'
                        : 'bg-red-100 text-red-800'
                    }`}
                  >
                    {district.is_active ? 'ใช้งาน' : 'ไม่ใช้งาน'}
                  </span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {formatDateOnly(district.created_at)}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                  <button
                    onClick={() => handleEdit(district)}
                    className="text-blue-600 hover:text-blue-900 mr-4"
                  >
                    แก้ไข
                  </button>
                  <button
                    onClick={() => handleDelete(district)}
                    className="text-red-600 hover:text-red-900"
                  >
                    ลบ
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>

        <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between p-4 border-t border-gray-200">
          <div className="flex items-center gap-2 text-sm text-gray-600">
            <span>แสดง</span>
            <select
              value={pageSize}
              onChange={(e) => handlePageSizeChange(Number(e.target.value))}
              className="border border-gray-300 rounded-md px-2 py-1 bg-white"
            >
              <option value={10}>10</option>
              <option value={20}>20</option>
              <option value={50}>50</option>
              <option value={100}>100</option>
            </select>
            <span>รายการต่อหน้า</span>
          </div>

          <div className="text-sm text-gray-600">
            หน้า {page} / {totalPages} (ทั้งหมด {total} รายการ)
          </div>

          <div className="flex items-center gap-2">
            <button
              onClick={() => handlePageChange(page - 1)}
              disabled={page <= 1}
              className="px-3 py-1 rounded-md border border-gray-300 text-gray-700 disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
            >
              ก่อนหน้า
            </button>
            <button
              onClick={() => handlePageChange(page + 1)}
              disabled={page >= totalPages}
              className="px-3 py-1 rounded-md border border-gray-300 text-gray-700 disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
            >
              ถัดไป
            </button>
          </div>
        </div>
      </div>

      <DistrictFormModal
        isOpen={isModalOpen}
        mode={modalMode}
        district={selectedDistrict}
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleSubmit}
      />
    </div>
  );
}
