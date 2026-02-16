import { apiClient } from './client';

export interface District {
  id: string;
  province_id: string;
  name: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface CreateDistrictDto {
  province_id: string;
  name: string;
  is_active: boolean;
}

export interface UpdateDistrictDto {
  province_id: string;
  name: string;
  is_active: boolean;
}

export interface ListDistrictsParams {
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

export interface ListDistrictsResponse {
  data: District[];
  paginate: {
    page: number;
    size: number;
    total: number;
  };
}

export const districtsApi = {
  async list(params?: ListDistrictsParams): Promise<ListDistrictsResponse> {
    const queryParams = new URLSearchParams();
    if (params?.page) queryParams.append('page', params.page.toString());
    if (params?.size) queryParams.append('size', params.size.toString());

    const query = queryParams.toString();
    const endpoint = `/system/districts/${query ? `?${query}` : ''}`;

    const response = await apiClient.get<District[]>(endpoint);
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

  async getById(id: string): Promise<District> {
    const response = await apiClient.get<District>(`/system/districts/${id}`);
    return response.data;
  },

  async create(data: CreateDistrictDto): Promise<void> {
    await apiClient.post('/system/districts/', data);
  },

  async update(id: string, data: UpdateDistrictDto): Promise<void> {
    await apiClient.patch(`/system/districts/${id}`, data);
  },

  async delete(id: string): Promise<void> {
    await apiClient.delete(`/system/districts/${id}`);
  },
};
