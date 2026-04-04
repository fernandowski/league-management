import {FlatList, View, StyleSheet, DimensionValue} from "react-native";
import {Text} from "react-native-paper";


interface TableProps<T> {
    data: T[];
    columns: ColumnDefinition<T>[]
}

export interface ColumnDefinition<T> {
    key: keyof T;
    title: string;
    width?: DimensionValue;
    render?: (item: T[keyof T]) => React.ReactNode;
}

export default function TableList<T>({ data, columns }: TableProps<T>) {
    const renderItem = ({ item }: { item: T }) => (
        <View style={styles.row}>
            {columns.map((column, index) => (
                <View key={index} style={[styles.cell, {width: column.width || 'auto'}]}>
                    {column.render ? column.render(item[column.key]) : <Text>{String(item[column.key])}</Text>}
                </View>
            ))}
        </View>
    );

    return (
        <View style={styles.container}>
            <View style={[styles.row, styles.header]}>
                {columns.map((column, index) => (
                    <View key={index} style={[styles.cell, { width: column.width || 'auto'}]}>
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
        backgroundColor: '#fff',
    },
    table: {
        marginTop: 8,
    },
    row: {
        flexDirection: 'row',
        borderBottomWidth: 1,
        borderBottomColor: '#ddd',
        paddingVertical: 8,
    },
    header: {
        backgroundColor: '#f0f0f0',
        borderBottomWidth: 2,
        borderBottomColor: '#bbb',
    },
    headerText: {
        fontWeight: 'bold',
        textAlign: "left"
    },
    cell: {
        paddingHorizontal: 8,
        justifyContent: 'center',
    },
});
