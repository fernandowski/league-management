import { StyleSheet, View } from 'react-native';
import { Button, Text } from 'react-native-paper';

import { useAppTheme } from '@/theme/theme';

interface PaginationProps {
  currentPage: number;
  totalItems: number;
  itemsPerPage: number;
  onPageChange: (newPage: number) => void;
  onItemsPerPageChange?: (newItemsPerPage: number) => void;
}

const Pagination: React.FC<PaginationProps> = ({
  currentPage,
  totalItems,
  itemsPerPage,
  onPageChange,
}) => {
  const totalPages = Math.ceil(totalItems / itemsPerPage);
  const from = currentPage * itemsPerPage + 1;
  const to = Math.min((currentPage + 1) * itemsPerPage, totalItems);
  const theme = useAppTheme();

  return (
    <View
      style={[
        styles.container,
        {
          backgroundColor: theme.colors.surface,
          borderColor: theme.colors.outline,
        },
      ]}
    >
      <Text style={[styles.info, { color: theme.colors.onSurfaceVariant }]}>
        Showing {from}–{to} of {totalItems} items
      </Text>
      <View style={styles.buttons}>
        <Button mode="text" disabled={currentPage === 0} onPress={() => onPageChange(currentPage - 1)}>
          Previous
        </Button>
        <Button
          mode="text"
          disabled={currentPage >= totalPages - 1}
          onPress={() => onPageChange(currentPage + 1)}
        >
          Next
        </Button>
      </View>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginTop: 16,
    paddingHorizontal: 16,
    paddingVertical: 12,
    borderWidth: 1,
    borderRadius: 18,
  },
  info: {
    fontSize: 14,
  },
  buttons: {
    flexDirection: 'row',
    gap: 8,
  },
});

export default Pagination;
