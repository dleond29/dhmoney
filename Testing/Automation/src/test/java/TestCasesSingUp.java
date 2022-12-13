import io.restassured.RestAssured;
import org.junit.jupiter.api.DisplayName;
import org.springframework.core.annotation.Order;
import org.testng.annotations.Test;

import static io.restassured.RestAssured.given;

public class TestCasesSingUp {

    private String baseUrl = "http://localhost:8080/api/users";

    @Test
    @Order(1)
    @DisplayName("singUpSuccessfully")
    public void SingUpSuccessfullyTestCase() {
        String register = "{\"name\":\"pepote\"," +
                "\"last_name\":\"mendez\"," +
                "\"dni\":\"1234567\"," +
                "\"phone\":\"3232323\"," +
                "\"email\":\"pepote@gmail.com\"," +
                "\"password\":\"pepoto\"}";

        RestAssured.baseURI = baseUrl;
        RestAssured
                .given().header("Content-Type", "application/json").body(register)
                .when().post("/")
                .then().assertThat().statusCode(200);
    }

    @Test
    @Order(2)
    @DisplayName("singUpRequiredFields")
    public void SingUpRequiredFieldsTestCase() {
        String register = "{\"name\":\"pepote\"," +
                "\"last_name\":\"mendez\"," +
                "\"email\":\"pepote@gmail.com\"," +
                "\"password\":\"pepoto\"}";

        RestAssured.baseURI = baseUrl;
        RestAssured
                .given().header("Content-Type", "application/json").body(register)
                .when().post("/")
                .then().assertThat().statusCode(400);
    }

    @Test
    @Order(3)
    @DisplayName("singUpEmailAlreadyRegistered")
    public void SingUpEmailAlreadyRegisteredTestCase() {
        String register = "{\"name\":\"pepote\"," +
                "\"last_name\":\"mendez\"," +
                "\"dni\":\"1234567\"," +
                "\"phone\":\"3232323\"," +
                "\"email\":\"pepote@gmail.com\"," +
                "\"password\":\"pepoto\"}";

        RestAssured.baseURI = baseUrl;
        RestAssured
                .given().header("Content-Type", "application/json").body(register)
                .when().post("/")
                .then().assertThat().statusCode(400);
    }


}
