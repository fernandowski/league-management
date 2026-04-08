import { DimensionValue, FlatList, StyleSheet, View } from 'react-native';
import { Text } from 'react-native-paper';

import { useAppTheme } from '@/theme/theme';

interface TableProps<T> {
  data: T[];
  columns: ColumnDefinition<T>[];
}

export interface ColumnDefinition<T> {
  key: keyof T;
  title: string;
  width?: DimensionValue;
  render?: (item: T[keyof T]) => React.ReactNode;
}

export default function TableList<T>({ data, columns }: TableProps<T>) {
  const theme = useAppTheme();

  const renderItem = ({ item }: { item: T }) => (
    <View style={[styles.row, { borderBottomColor: theme.colors.outline }]}>
      {columns.map((column, index) => (
        <View key={index} style={[styles.cell, { width: column.width || 'auto' }]}>
          {column.render ? column.render(item[column.key]) : <Text>{String(item[column.key])}</Text>}
        </View>
      ))}
    </View>
  );

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
      <View
        style={[
          styles.row,
          styles.header,
          {
            backgroundColor: theme.colors.surfaceVariant,
            borderBottomColor: theme.colors.outlineVariant,
          },
        ]}
      >
        {columns.map((column, index) => (
          <View key={index} style={[styles.cell, { width: column.width || 'auto' }]}>
            <Text style={styles.headerText}>{column.title}</Text>
          </View>
        ))}
      </View>

      <FlatList
        data={data}
        renderItem={renderItem}
        keyExtractor={(_, index) => index.toString()}
        contentContainerStyle={styles.table}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    padding: 16,
    borderWidth: 1,
    borderRadius: 24,
  },
  table: {
    marginTop: 8,
  },
  row: {
    flexDirection: 'row',
    borderBottomWidth: 1,
    paddingVertical: 8,
  },
  header: {
    borderBottomWidth: 2,
    borderTopLeftRadius: 16,
    borderTopRightRadius: 16,
  },
  headerText: {
    fontWeight: 'bold',
    textAlign: 'left',
  },
  cell: {
    paddingHorizontal: 8,
    justifyContent: 'center',
  },
});
