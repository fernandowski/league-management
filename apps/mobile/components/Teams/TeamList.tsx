import {Animated, View} from "react-native";
import ScrollView = Animated.ScrollView;
import {Card, Text} from "react-native-paper";
import {Team} from "@/components/Teams/TeamOverview";

interface props {
    data: Team[]
}
export default function TeamList (props: props) {
    return (
        <View style={{flex: 1}}>
            <ScrollView>
                {
                    props.data.map((item) => (
                        <Card style={{marginBottom: 8, marginLeft: 1, marginRight: 1}} key={item.id}>
                            <Card.Content>
                                <Text>{item.name}</Text>
                            </Card.Content>
                        </Card>
                    ))
                }
            </ScrollView>
        </View>
    )
}
