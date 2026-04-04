import {Dimensions, StyleSheet, Text, TouchableOpacity, useWindowDimensions, View} from "react-native";
import {DrawerNavigationOptions, DrawerNavigationProp} from "@react-navigation/drawer";
import {ParamListBase, RouteProp} from "@react-navigation/native";
import {Select} from "@/components/Select/Select";
import {useOrganizationStore} from "@/stores/organizationStore";
import {useEffect} from "react";
import {useNavigation} from "expo-router";
import Ionicons from "@expo/vector-icons/Ionicons";


interface HeaderProps {
    navigation: DrawerNavigationProp<ParamListBase>
    route: RouteProp<ParamListBase>
    options: DrawerNavigationOptions
}

export function Header() {
    const navigation: DrawerNavigationProp<any> = useNavigation();
    const dimensions = useWindowDimensions();
    const isLargeScreen = dimensions.width >= 768;

    const {organizations, loading, error, fetchOrganizations, setOrganization, organization} = useOrganizationStore();
    const onSelectChange = (value: string) => {
        setOrganization(value);
    }

    useEffect(() => {
        fetchOrganizations();
    }, [])

    return (
        <View style={[styles.view]}>
            { !isLargeScreen && (
                <TouchableOpacity onPress={() => navigation.openDrawer()}>
                    <Ionicons name="menu" size={24} color="black" />
                </TouchableOpacity>
            )}

            <View style={[styles.select]}>
                <Text>Organizations: </Text>
                <Select onChange={onSelectChange} selected={organization} data={organizations.map(organizations => ({label: organizations.name, value: organizations.id}))}/>
            </View>
        </View>
    )
}

const {height: screenHeight} = Dimensions.get('window');

const styles = StyleSheet.create({
    view: {
        height: screenHeight * .05,
    },
    select: {
        flex: 1,
        flexDirection: "row",
        gap: 16,
        marginTop: 18
    }
});
