import { apiClient } from './client';

export interface Status {
  id: string;
  name_th: string;
  name_en: string;
  is_active: boolean;
  created_at: string;
}

export interface StatusListParams {
  page?: number;
  size?: number;
}

export interface StatusCreateInput {
  name_th: string;
  name_en: string;
  is_active: boolean;
}

export interface StatusUpdateInput {
  name_th: string;
  name_en: string;
  is_active: boolean;
}

interface BackendPaginate {
  page?: number;
  size?: number;
  total?: number;
  Page?: number;
  Size?: number;
  Total?: number;
}

export interface StatusListResponse {
  data: Status[];
  paginate: {
    page: number;
    size: number;
    total: number;
  };
}

export const statusesApi = {
  async list(params?: StatusListParams): Promise<StatusListResponse> {
    const queryParams = new URLSearchParams();
    if (params?.page) queryParams.append('page', params.page.toString());
    if (params?.size) queryParams.append('size', params.size.toString());

    const query = queryParams.toString();
    const endpoint = `/system/statuses${query ? `?${query}` : ''}`;

    const response = await apiClient.get<Status[]>(endpoint);
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

  async getById(id: string): Promise<Status> {
    const response = await apiClient.get<Status>(`/system/statuses/${id}`);
    return response.data;
  },

  async create(data: StatusCreateInput): Promise<void> {
    await apiClient.post('/system/statuses', data);
  },

  async update(id: string, data: StatusUpdateInput): Promise<void> {
    await apiClient.patch(`/system/statuses/${id}`, data);
  },

  async delete(id: string): Promise<void> {
    await apiClient.delete(`/system/statuses/${id}`);
  },
};
