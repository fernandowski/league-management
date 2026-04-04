import {Animated, StyleSheet, View} from "react-native";
import {useOrganizationStore} from "@/stores/organizationStore";
import {Card, Text} from "react-native-paper";
import ScrollView = Animated.ScrollView;

export function OrganizationTable() {
    const {organizations} = useOrganizationStore();
    return (
        <View style={[styles.outerContainer]}>
            <ScrollView style={[styles.viewContainer]}>
                {organizations.map(organization =>
                    <Card style={{marginTop: 8,}} key={organization.id}>
                        <Card.Content>
                            <Text>{organization.name}</Text>
                        </Card.Content>
                    </Card>)}
            </ScrollView>
        </View>
    )
}


const styles = StyleSheet.create({
    outerContainer: {
        flex: 1,
        justifyContent: 'center',
        alignItems: 'center'
    },
    viewContainer: {
        flex: 1,
        width: '80%',
        paddingLeft: 8,
        paddingRight: 8
    }
})
