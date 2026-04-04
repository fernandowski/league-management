import {GestureHandlerRootView} from "react-native-gesture-handler";
import {Drawer} from "expo-router/drawer";
import {useWindowDimensions} from "react-native";

export default function Layout() {
    const dimensions = useWindowDimensions();
    const isLargeScreen = dimensions.width >= 768;

    return (
        <GestureHandlerRootView style={{flex: 1}}>
            <Drawer
                screenOptions={{
                    drawerType: isLargeScreen ? 'permanent' : undefined,
                    header: () => null,
                    drawerPosition: 'left',
                    unmountOnBlur: true,
                }}
            >
                <Drawer.Screen
                    name="index"
                    options={{
                        drawerLabel: 'Organizations Overview',
                        title: 'Organizations Overview',
                    }}
                />
                <Drawer.Screen
                    name="leagues/[id]/index"
                    options={{
                        drawerLabel: 'Leagues',
                        title: 'Leagues',
                    }}
                />
                <Drawer.Screen
                    name="teams/index"
                    options={{
                        drawerLabel: 'Teams',
                        title: 'Teams',
                    }}
                />
                <Drawer.Screen
                    name="seasons/index"
                    options={{
                        drawerLabel: 'Seasons',
                        title: 'Seasons',
                    }}
                />

                <Drawer.Screen
                    name="seasons/[id]"
                    options={{
                        drawerItemStyle: {display: 'none'},
                        drawerLabel: () => null,
                    }}
                />
            </Drawer>
        </GestureHandlerRootView>
    );
}
