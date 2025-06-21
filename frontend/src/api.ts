import {useFetch} from "@vueuse/core";

const baseUrl = '/api';

export async function useChannels() {
  const { data, error, isFetching, execute, abort } = useFetch(`${baseUrl}/channels`);
  return useFetch(`${baseUrl}/channels`);

  return {
    data,
    error,
    isFetching,
    execute,
    abort,
  };
}
