import { apiClient } from './client';

export interface SubDistrict {
  id: string;
  district_id: string;
  name: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface CreateSubDistrictDto {
  district_id: string;
  name: string;
  is_active: boolean;
}

export interface UpdateSubDistrictDto {
  district_id: string;
  name: string;
  is_active: boolean;
}

export interface ListSubDistrictsParams {
  page?: number;
  size?: number;
}

interface BackendPaginate {
  page?: number;
  size?: number;
  total?: number;
  Page?: number;
  Size?: number;
  Total?: number;
}

export interface ListSubDistrictsResponse {
  data: SubDistrict[];
  paginate: {
    page: number;
    size: number;
    total: number;
  };
}

export const subDistrictsApi = {
  async list(params?: ListSubDistrictsParams): Promise<ListSubDistrictsResponse> {
    const queryParams = new URLSearchParams();
    if (params?.page) queryParams.append('page', params.page.toString());
    if (params?.size) queryParams.append('size', params.size.toString());

    const query = queryParams.toString();
    const endpoint = `/system/sub_districts/${query ? `?${query}` : ''}`;

    const response = await apiClient.get<SubDistrict[]>(endpoint);
    const rawPaginate = (response.paginate ?? {}) as BackendPaginate;

    return {
      data: response.data,
      paginate: {
        page: rawPaginate.page ?? rawPaginate.Page ?? 1,
        size: rawPaginate.size ?? rawPaginate.Size ?? params?.size ?? 10,
        total: rawPaginate.total ?? rawPaginate.Total ?? response.data.length,
      },
    };
  },

  async getById(id: string): Promise<SubDistrict> {
    const response = await apiClient.get<SubDistrict>(`/system/sub_districts/${id}`);
    return response.data;
  },

  async create(data: CreateSubDistrictDto): Promise<void> {
    await apiClient.post('/system/sub_districts/', data);
  },

  async update(id: string, data: UpdateSubDistrictDto): Promise<void> {
    await apiClient.patch(`/system/sub_districts/${id}`, data);
  },

  async delete(id: string): Promise<void> {
    await apiClient.delete(`/system/sub_districts/${id}`);
  },
};
