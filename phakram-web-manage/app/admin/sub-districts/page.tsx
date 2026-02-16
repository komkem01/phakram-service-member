"use client";

import { useState, useEffect, useCallback, useRef } from "react";
import { SubDistrict, subDistrictsApi } from "@/lib/api/sub-districts";
import { districtsApi } from "@/lib/api/districts";
import { formatDateOnly } from "@/lib/utils/date";
import Loading from "@/components/admin/Loading";
import SubDistrictFormModal from "@/components/admin/SubDistrictFormModal";
import { useModal } from "@/contexts/ModalContext";

export default function SubDistrictsPage() {
  const [subDistricts, setSubDistricts] = useState<SubDistrict[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [total, setTotal] = useState(0);
  const [districtNameMap, setDistrictNameMap] = useState<Record<string, string>>({});
  const districtNameMapRef = useRef<Record<string, string>>({});
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [modalMode, setModalMode] = useState<'create' | 'edit'>('create');
  const [selectedSubDistrict, setSelectedSubDistrict] = useState<SubDistrict | null>(null);

  const { showSuccess, showError, showWarning } = useModal();

  const totalPages = Math.max(1, Math.ceil(total / pageSize));

  useEffect(() => {
    districtNameMapRef.current = districtNameMap;
  }, [districtNameMap]);

  const hydrateMissingDistrictNames = useCallback(async (items: SubDistrict[]) => {
    const neededIds = Array.from(new Set(items.map((item) => item.district_id).filter(Boolean)));
    const missingIds = neededIds.filter((id) => !districtNameMapRef.current[id]);
    if (missingIds.length === 0) return;

    const results = await Promise.allSettled(
      missingIds.map(async (id) => {
        const district = await districtsApi.getById(id);
        return { id, name: district.name };
      })
    );

    const fallbackMap = results.reduce<Record<string, string>>((acc, result) => {
      if (result.status === 'fulfilled') {
        acc[result.value.id] = result.value.name;
      }
      return acc;
    }, {});

    if (Object.keys(fallbackMap).length > 0) {
      setDistrictNameMap((prev) => ({ ...prev, ...fallbackMap }));
    }
  }, []);

  const fetchDistricts = useCallback(async () => {
    try {
      const response = await districtsApi.list({ page: 1, size: 1000 });
      const mapping = response.data.reduce<Record<string, string>>((acc, district) => {
        acc[district.id] = district.name;
        return acc;
      }, {});
      setDistrictNameMap(mapping);
    } catch (error) {
      console.error('Failed to load districts:', error);
    }
  }, []);

  const fetchSubDistricts = useCallback(async (currentPage: number, currentSize: number) => {
    try {
      setIsLoading(true);
      const response = await subDistrictsApi.list({ page: currentPage, size: currentSize });
      setSubDistricts(response.data);
      setTotal(response.paginate.total);
      setPage(response.paginate.page);
      await hydrateMissingDistrictNames(response.data);
    } catch {
      showError('ไม่สามารถโหลดข้อมูลตำบลได้');
    } finally {
      setIsLoading(false);
    }
  }, [hydrateMissingDistrictNames, showError]);

  useEffect(() => {
    fetchSubDistricts(page, pageSize);
  }, [fetchSubDistricts, page, pageSize]);

  useEffect(() => {
    fetchDistricts();
  }, [fetchDistricts]);

  const handleCreate = () => {
    setModalMode('create');
    setSelectedSubDistrict(null);
    setIsModalOpen(true);
  };

  const handleEdit = (subDistrict: SubDistrict) => {
    setModalMode('edit');
    setSelectedSubDistrict(subDistrict);
    setIsModalOpen(true);
  };

  const handleDelete = (subDistrict: SubDistrict) => {
    showWarning(
      `ต้องการลบตำบล "${subDistrict.name}" หรือไม่?`,
      async () => {
        try {
          await subDistrictsApi.delete(subDistrict.id);
          showSuccess('ลบตำบลสำเร็จ');
          fetchSubDistricts(page, pageSize);
        } catch {
          showError('ไม่สามารถลบตำบลได้');
        }
      },
      'ยืนยันการลบ'
    );
  };

  const handleSubmit = async (data: { district_id: string; name: string; is_active: boolean }) => {
    try {
      if (modalMode === 'create') {
        await subDistrictsApi.create(data);
        showSuccess('เพิ่มตำบลสำเร็จ');
      } else if (selectedSubDistrict) {
        await subDistrictsApi.update(selectedSubDistrict.id, data);
        showSuccess('แก้ไขตำบลสำเร็จ');
      }
      setIsModalOpen(false);
      fetchSubDistricts(page, pageSize);
    } catch {
      showError(`ไม่สามารถ${modalMode === 'create' ? 'เพิ่ม' : 'แก้ไข'}ตำบลได้`);
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
        <h1 className="text-3xl font-bold text-gray-800">จัดการตำบล</h1>
        <button
          onClick={handleCreate}
          className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded-lg font-medium transition-colors"
        >
          + เพิ่มตำบล
        </button>
      </div>

      <div className="bg-white rounded-lg shadow overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                ชื่อตำบล
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                อำเภอ
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
            {subDistricts.map((subDistrict) => (
              <tr key={subDistrict.id}>
                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                  {subDistrict.name}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-700">
                  {districtNameMap[subDistrict.district_id] ?? '-'}
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span
                    className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                      subDistrict.is_active
                        ? 'bg-green-100 text-green-800'
                        : 'bg-red-100 text-red-800'
                    }`}
                  >
                    {subDistrict.is_active ? 'ใช้งาน' : 'ไม่ใช้งาน'}
                  </span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {formatDateOnly(subDistrict.created_at)}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                  <button
                    onClick={() => handleEdit(subDistrict)}
                    className="text-blue-600 hover:text-blue-900 mr-4"
                  >
                    แก้ไข
                  </button>
                  <button
                    onClick={() => handleDelete(subDistrict)}
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

      <SubDistrictFormModal
        isOpen={isModalOpen}
        mode={modalMode}
        subDistrict={selectedSubDistrict}
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleSubmit}
      />
    </div>
  );
}
